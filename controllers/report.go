package controllers

import (
	"net/http"
	"perpustakaan-golang/models" // Import your models package
	"perpustakaan-golang/types"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController struct{} 

// BorrowedBookReports retrieves data for borrowed book reports.
// @Summary Borrowed Book Reports
// @Description Get borrowed book reports within a date range.
// @Tags Reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} types.StandardResponse(data=[]models.BorrowedBookReportSwagger)
// @Failure 400 {object} models.ErrorResponse 
// @Failure 500 {object} models.ErrorResponse
// @Router /reports/borrowed_books [get]
func (rc ReportController) BorrowedBookReports(c *gin.Context) {
	var reportRequest struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	}
	if err := c.ShouldBind(&reportRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates or use default values
	startDate, err := time.Parse("2006-01-02", reportRequest.StartDate)
	if err != nil {
		startDate = time.Time{} 
	}
	endDate, err := time.Parse("2006-01-02", reportRequest.EndDate)
	if err != nil {
		endDate = time.Now()
	}


	report, err := models.GetBorrowedBookReports(startDate, endDate) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	reportSwagger := make([]models.BorrowedBookReportSwagger, len(report))
	for i, r := range report {
		reportSwagger[i] = models.BorrowedBookReportSwagger{
			IDBuku:           r.IDBuku,
			Judul:            r.Judul,
			Pengarang:        r.Pengarang,
			IDAnggota:        r.IDAnggota,
			Nama:             r.Nama,
			TanggalPinjam:    r.TanggalPinjam,
			TanggalJatuhTempo: r.TanggalJatuhTempo,
			TanggalKembali:   r.TanggalKembali.Time.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Get borrowed book reports successfully",
		Data:    reportSwagger,
	})
}

// ReturnedBookReports retrieves data for returned book reports.
// @Summary Returned Book Reports
// @Description Get returned book reports within a date range.
// @Tags Reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} types.StandardResponse(data=[]models.ReturnedBookReportSwagger)
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reports/returned_books [get]
func (rc ReportController) ReturnedBookReports(c *gin.Context) {
	var reportRequest struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	}
	if err := c.ShouldBind(&reportRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates or use default values
	startDate, err := time.Parse("2006-01-02", reportRequest.StartDate)
	if err != nil {
		startDate = time.Time{} 
	}
	endDate, err := time.Parse("2006-01-02", reportRequest.EndDate)
	if err != nil {
		endDate = time.Now()
	}

	report, err := models.GetReturnedBookReports(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	reportSwagger := make([]models.ReturnedBookReportSwagger, len(report))
    for i, r := range report {
        reportSwagger[i] = models.ReturnedBookReportSwagger{
            IDBuku:           r.IDBuku,
            Judul:            r.Judul,
            Pengarang:        r.Pengarang,
            IDAnggota:        r.IDAnggota,
            Nama:             r.Nama,
            TanggalPinjam:    r.TanggalPinjam,
            TanggalJatuhTempo: r.TanggalJatuhTempo,
            TanggalKembali:   r.TanggalKembali,
            Denda:            r.Denda.Int64,
        }
    }

	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Get returned book reports successfully",
		Data:    reportSwagger,
	})
}