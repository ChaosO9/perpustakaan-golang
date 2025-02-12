package controllers

import (
	"net/http"
	"strconv"

	"perpustakaan-golang/forms"
	"perpustakaan-golang/models"
	"perpustakaan-golang/types"
	"perpustakaan-golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateBookFormSwagger struct {
    Judul          string `json:"judul" binding:"required,min=3,max=255"`
    Pengarang      string `json:"pengarang" binding:"required,min=3,max=255"`
    Penerbit       string `json:"penerbit" binding:"required,min=3,max=255"`
    ISBN           string `json:"isbn" binding:"required,min=10,max=17"`
    TahunTerbit    int64  `json:"tahun_terbit" binding:"required,numeric"`
    Kategori       string `json:"kategori" binding:"required,min=3,max=255"`
    Deskripsi      string `json:"deskripsi" binding:"required,min=3,max=1000"`
    JumlahEksemplar int64  `json:"jumlah_eksemplar" binding:"required,numeric,min=1"`
}

type BookAll struct {
	Books []models.BookSwagger `json:"books"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

type BookController struct{}

var bookModel = new(models.BookModel)
var bookForm = new(forms.BookForm)

// Create adds a new book.
// @Summary Create Book
// @Description Add a new book
// @Tags Books
// @Accept multipart/form-data 
// @Produce json
// @Param request body controllers.CreateBookFormSwagger true "Create Book Request"
// @Success 200 {object} types.StandardResponse{data=models.BookCreatedResponse}
// @Failure 406 {object} models.ErrorResponse 
// @Router /books [post]
func (ctrl BookController) Create(c *gin.Context) {
	var form forms.CreateBookForm
	if err := c.ShouldBind(&form); err != nil {
		message := bookForm.Create(err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	uploadDir := "./uploads/book/"
	filename := "book-cover-placeholder.png"  
	if form.Foto != nil { 
        var err error 
		filename, err = utils.ImageUpload(c, "foto", uploadDir)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
	}
	form.FotoFileName = filename

    filename, err := utils.ImageUpload(c, "foto", uploadDir) // "foto" is the form field name
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	form.FotoFileName = filename

	id, err := bookModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Book created successfully",
		Data:   models.BookCreatedResponse{
			Message: "Book created",
			ID:      id,
		},
	})
}

// All retrieves all books with pagination.
// @Summary Get All Books
// @Description Retrieve all books with pagination
// @Tags Books
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} types.StandardResponse{data=BookAll}
// @Router /books [get]
func (ctrl BookController) All(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 10 // Default limit

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	books, total, err := bookModel.All(page, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	booksSwagger := models.ToBooksSwagger(books)

	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Books retrieved successfully",
		Data: BookAll{
			Books: booksSwagger,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

// One retrieves a single book by ID.
// @Summary Get Book by ID
// @Description Retrieve a book by its ID
// @Tags Books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} types.StandardResponse{data=models.BookSwagger}
// @Failure 404 {object} models.ErrorResponse 
// @Router /books/{id} [get]
func (ctrl BookController) One(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid book ID"})
		return
	}
	book, err := bookModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	bookSwagger := book.ToBookSwagger()

	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Book retrieved successfully",
		Data:    bookSwagger,
	})
}

// Update modifies an existing book.
// @Summary Update Book
// @Description Update an existing book
// @Tags Books
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Book ID"
// @Param request body controllers.CreateBookFormSwagger true "Update Book Request"
// @Success 200 {object} types.StandardResponse{data=nil}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 406 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [put]
func (ctrl BookController) Update(c *gin.Context) {
    uploadDir := "./uploads/book/"
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid book ID"})
        return
    }

    var form forms.CreateBookForm
    if err := c.ShouldBind(&form); err != nil {
        message := bookForm.Update(err)
        c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
        return
    }

    existingBook, err := bookModel.One(id)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get book"})
        return
    }
    oldFilename := existingBook.Foto.String
	form.FotoFileName = oldFilename

    // Only update the image if a new file was uploaded
    if form.Foto != nil {
        newFilename, err := utils.ImageUpdate(c, "foto", uploadDir, oldFilename)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        form.FotoFileName = newFilename
    }


    err = bookModel.Update(id, form)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, types.StandardResponse{
        Code:    http.StatusOK,
        Message: "Book updated successfully",
        Data:    nil,
    })
}

// Delete removes a book.
// @Summary Delete Book
// @Description Delete a book by its ID
// @Tags Books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} types.StandardResponse{data=nil}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [delete]
func (ctrl BookController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid book ID"})
		return
	}

	existingBook, err := bookModel.One(id)  
	if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get book"})
        return
    }
    filename := existingBook.Foto.String

	uploadDir := "./uploads/book/"
	if err := utils.ImageDelete(uploadDir, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bookModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Book deleted successfully",
		Data:    nil,
	})
}
