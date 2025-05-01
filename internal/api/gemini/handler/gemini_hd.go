package searchHandler

import ( 
	"ProjectGolang/pkg/log"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/net/context"
	"time"
	"github.com/gofiber/fiber/v2"
)

// Lengkapi function Handler untuk mengunggah dan memproses gambar
func (h *SearchHandler) AnalyzeMedicineImage(ctx *fiber.Ctx) error { 
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("image")
	if err != nil { 
		h.log.WithFields(log.Fields{ 
			"error" : err.Error(),
		}).Error("Failed to get image ")
		return err
	}
	
	medicineResponse, err := h.searchService.AnalyzeMedicineImage(c, file) 
	if err != nil { 
		h.log.WithFields(log.Fields{ 
			"error": err.Error(),
		}).Error("Failed to Analyze Medicine Image")
	return err
	}
	select { 
	case <-c.Done(): 
	return ctx.Status(fiber.StatusRequestTimeout).JSON(utils.StatusMessage(fiber.StatusRequestTimeout))
	default: 
	return ctx.Status(fiber.StatusOK).JSON(medicineResponse)
	}	
}
