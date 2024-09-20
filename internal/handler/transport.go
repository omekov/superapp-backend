package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/omekov/dubaicarkzv2/internal/usecase"
)

type homeHandler struct {
	useCase usecase.UseCase
}

func (h homeHandler) handlerTransport(w http.ResponseWriter, r *http.Request) error {
	mark := strings.ToUpper(r.URL.Query().Get("mark"))
	model := strings.ToUpper(r.URL.Query().Get("model"))
	volumeQuery := r.URL.Query().Get("volume")
	if mark != "" && model == "" && volumeQuery == "" {
		models, err := h.useCase.GetModels(r.Context(), mark)
		if err != nil {
			return err
		}

		modelsByte, err := json.Marshal(models)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		w.Write(modelsByte)
		return nil
	}

	if mark != "" && model != "" && volumeQuery == "" {
		volumes, err := h.useCase.GetVolumes(r.Context(), mark, model)
		if err != nil {
			return err
		}

		volumesByte, err := json.Marshal(volumes)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		w.Write(volumesByte)
		return nil
	}

	if mark != "" && model != "" && volumeQuery != "" {
		volume, err := strconv.Atoi(volumeQuery)
		if err != nil {
			return fmt.Errorf("volume -> %v", err)
		}
		fmt.Println(volume)
		specifications, err := h.useCase.GetSpecifications(r.Context(), mark, model, volume)
		if err != nil {
			return err
		}

		specificationsByte, err := json.Marshal(specifications)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		w.Write(specificationsByte)
		return nil
	}

	marks, err := h.useCase.GetMarks(r.Context())
	if err != nil {
		return err
	}

	marksByte, err := json.Marshal(marks)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marksByte)
	return nil
}

func (h homeHandler) handlerAssessment(w http.ResponseWriter, r *http.Request) error {
	return nil
}
