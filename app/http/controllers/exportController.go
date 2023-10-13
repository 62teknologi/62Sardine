package controllers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/xuri/excelize/v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/62teknologi/62sardine/config"

	"github.com/gin-gonic/gin"
)

type ExportController struct {
}

type ExportRequest struct {
	FileName  string    `json:"file_name"`
	Template  string    `json:"template"`
	SheetName string    `json:"sheet_name"`
	Headings  []Heading `json:"headings"`
	Data      []map[string]any
}

type Heading struct {
	Field     string `json:"field"`
	FieldName string `json:"field_name"`
}

func (ctrl *ExportController) Export(ctx *gin.Context) {
	exportFolder, _ := config.ReadConfig("filesystems.export_folder")

	exportTo := ctx.Query("export_to")
	if !(exportTo == "xlsx" || exportTo == "csv" || exportTo == "pdf") {
		ctrl.ResErr(ctx, errors.New("invalid export type, only support either xlsx or csv or pdf"))
		return
	}

	var exportRequest ExportRequest

	if err := ctx.BindJSON(&exportRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while get request data: " + err.Error()})
		return
	}

	if exportTo == "csv" {
		// Create a CSV writer.
		csvFileName := exportFolder + "/" + exportRequest.FileName + ".csv"
		file, err := os.Create(csvFileName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while create a csv writer"})
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write headings to the CSV file.
		var headings []string
		for _, d := range exportRequest.Headings {
			headings = append(headings, d.FieldName)
		}
		if err := writer.Write(headings); err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while write headings to the csv file"})
			return
		}

		// Write data to the CSV file.
		for _, exportData := range exportRequest.Data {
			var row []string
			for _, d := range exportRequest.Headings {
				row = append(row, fmt.Sprintf("%v", exportData[d.Field]))
			}
			if err := writer.Write(row); err != nil {
				fmt.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while write data to the csv file"})
				return
			}
		}
	} else if exportTo == "xlsx" {
		xlsx := excelize.NewFile()
		defer func() {
			if err := xlsx.Close(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while create new xlsx file"})
			}
		}()

		// Create a new sheet.
		index, err := xlsx.NewSheet(exportRequest.SheetName)
		if err != nil {
			fmt.Println(err)
			return
		}

		rowIndex := 1
		for i, d := range exportRequest.Headings {
			rowName := string(rune('A' + i))
			cellName := rowName + strconv.Itoa(rowIndex)
			xlsx.SetCellValue(exportRequest.SheetName, cellName, d.FieldName)
			style, _ := xlsx.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
			xlsx.SetCellStyle(exportRequest.SheetName, cellName, cellName, style)

			for exportIteration, exportData := range exportRequest.Data {
				cellName := rowName + strconv.Itoa(rowIndex+exportIteration+1)
				// Set value of a cell.
				xlsx.SetCellValue(exportRequest.SheetName, cellName, exportData[d.Field])
			}
		}

		// Set active sheet of the workbook.
		xlsx.SetActiveSheet(index)

		// Save spreadsheet by the given path.
		if err := xlsx.SaveAs(exportFolder + "/" + exportRequest.FileName + ".xlsx"); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while save file"})
			return
		}
	}

	// Set the file path to the file you want to serve for download
	filePath := exportFolder + "/" + exportRequest.FileName + "." + exportTo

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "error while open file"})
		return
	}

	defer file.Close()

	// Get the file information
	fileInfo, err := file.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while get file info"})
		return
	}

	// Set the response headers
	ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", fmt.Sprint(fileInfo.Size()))

	// Stream the file content as the response body
	http.ServeContent(ctx.Writer, ctx.Request, fileInfo.Name(), fileInfo.ModTime(), file)
}

type ExportPDFRequest struct {
	OutputName     string         `json:"output_name"`
	Template       string         `json:"template"`
	Data           map[string]any `json:"data"`
	ParsedTemplate string
}

func (r *ExportPDFRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	r.ParsedTemplate = buf.String()
	return nil
}

func (r *ExportPDFRequest) GeneratePDF(pdfPath string, args []string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("cloneTemplate/", 0777)
		if errDir != nil {
			log.Fatal(errDir)
		}
	}
	err1 := os.WriteFile("cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.ParsedTemplate), 0644)
	if err1 != nil {
		panic(err1)
	}

	f, err := os.Open("cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	//We need install this to OS
	//sudo apt install wkhtmltopdf - for Ubuntu
	//brew install Caskroom/cask/wkhtmltopdf - For Mac
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	for _, arg := range args {
		switch arg {
		case "low-quality":
			pdfg.LowQuality.Set(true)
		case "no-pdf-compression":
			pdfg.NoPdfCompression.Set(true)
		case "grayscale":
			pdfg.Grayscale.Set(true)
			// Add other arguments as needed
		}
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {

		}
	}(dir + "/cloneTemplate")

	return true, nil
}

func (ctrl *ExportController) ExportPDF(ctx *gin.Context) {
	var request ExportPDFRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while get request data: " + err.Error()})
		return
	}

	//html template path
	templatePath := "./public/" + request.Template

	//path for download pdf
	outputPath := "./storage/" + request.OutputName + ".pdf"

	if err := request.ParseTemplate(templatePath, request.Data); err == nil {

		// TODO : add more compression options
		// Generate PDF with custom arguments
		args := []string{"no-pdf-compression"}

		// Generate PDF
		ok, _ := request.GeneratePDF(outputPath, args)
		fmt.Println(ok, "pdf generated successfully")

		file, err := os.Open(outputPath)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "error while open file"})
			return
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		fileInfo, err := file.Stat()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while get file info"})
			return
		}

		ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Length", fmt.Sprint(fileInfo.Size()))

		http.ServeContent(ctx.Writer, ctx.Request, fileInfo.Name(), fileInfo.ModTime(), file)
	} else {
		fmt.Println(err)
	}
}

func (ctrl *ExportController) ResErr(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}
