package controllers

import (
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/notefan-golang/helpers/chanh"
	"github.com/notefan-golang/helpers/combh"
	"github.com/notefan-golang/helpers/decodeh"
	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/nullh"
	"github.com/notefan-golang/helpers/rwh"
	"github.com/notefan-golang/helpers/validationh"
	"github.com/notefan-golang/models/requests/file_reqs"
	pc_reqs "github.com/notefan-golang/models/requests/page_content_reqs"
	"github.com/notefan-golang/models/responses"
	"github.com/notefan-golang/policies"
	"github.com/notefan-golang/services"
)

type PageContentController struct {
	service *services.PageContentService
	policy  *policies.PageContentPolicy
}

func NewPageContentController(
	service *services.PageContentService,
	policy *policies.PageContentPolicy,
) *PageContentController {
	return &PageContentController{
		service: service,
		policy:  policy,
	}
}

/*
 * ----------------------------------------------------------------
 * Controller utility methods ⬇
 * ----------------------------------------------------------------
 */

//

/*
 * ----------------------------------------------------------------
 * Controller handler func methods ⬇
 * ----------------------------------------------------------------
 */

// Get gets PageContents by request form data
func (controller PageContentController) Get(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[pc_reqs.GetByPage](r.Form)
	errorh.Panic(err)
	if orderBys, ok := r.Form["order_bys"]; ok {
		input.OrderBys = orderBys
	}

	err = controller.policy.Get(r.Context(), input)
	errorh.Panic(err)

	pageContentPagination, err := controller.service.GetBypage(r.Context(), input)
	errorh.Panic(err)

	pageContentPagination.SetPage(input.PerPage, input.Page, nullh.NullInt())
	pageContentPagination.SetURL(r.URL)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve contents of page").
			Body("contents", pageContentPagination),
	)
}

// Find finds a PageContent by request form data
func (controller PageContentController) Find(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[pc_reqs.Action](r.Form)
	errorh.Panic(err)

	err = controller.policy.Find(r.Context(), input)
	errorh.Panic(err)

	pageRes, err := controller.service.Find(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully retrieve page").
			Body("page", pageRes),
	)
}

// Create creates a PageContent from given request form data
func (controller PageContentController) Create(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[pc_reqs.Create](r.Form)
	errorh.Panic(err)

	// Load Medias from Request if exists
	wg, mutex, errChan := new(sync.WaitGroup), sync.Mutex{}, chanh.Make[error](nil, 1)
	defer close(errChan)
	if mediaFHs, ok := r.MultipartForm.File["medias"]; ok {
		for _, mediaFH := range mediaFHs {
			wg.Add(1)
			go func(mediaFH *multipart.FileHeader) {
				defer wg.Done()
				errChanVal := chanh.GetValAndKeep(errChan)
				if errChanVal != nil {
					return
				}
				mediaFileReq, errChanVal := file_reqs.NewFromFH(mediaFH)
				if errChanVal != nil {
					chanh.ReplaceVal(errChan, errChanVal)
					return
				}
				mutex.Lock()
				defer mutex.Unlock()
				input.Medias = append(input.Medias, mediaFileReq)
			}(mediaFH)
		}
	}
	wg.Wait()
	errorh.Panic(<-errChan)

	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Create(r.Context(), input)
	errorh.Panic(err)

	pageRes, err := controller.service.Create(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully create page").
			Body("page", pageRes),
	)
}

// Update updates a PageContent by request form data
func (controller PageContentController) Update(w http.ResponseWriter, r *http.Request) {
	input, err := decodeh.FormData[pc_reqs.Update](r.Form)
	errorh.Panic(err)

	// Load Medias from Request if exists
	wg, mutex, errChan := new(sync.WaitGroup), sync.Mutex{}, chanh.Make[error](nil, 1)
	defer close(errChan)
	if mediaFHs, ok := r.MultipartForm.File["medias"]; ok {
		for _, mediaFH := range mediaFHs {
			wg.Add(1)
			go func(mediaFH *multipart.FileHeader) {
				defer wg.Done()
				errChanVal := chanh.GetValAndKeep(errChan)
				if errChanVal != nil {
					return
				}
				mediaFileReq, errChanVal := file_reqs.NewFromFH(mediaFH)
				if errChanVal != nil {
					chanh.ReplaceVal(errChan, errChanVal)
					return
				}
				mutex.Lock()
				defer mutex.Unlock()
				input.Medias = append(input.Medias, mediaFileReq)
			}(mediaFH)
		}
	}
	wg.Wait()
	errorh.Panic(<-errChan)

	err = validationh.ValidateStruct(input)
	errorh.Panic(err)

	err = controller.policy.Update(r.Context(), input)
	errorh.Panic(err)

	pageContentRes, err := controller.service.Update(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w,
		responses.NewResponse().
			Code(http.StatusOK).
			Success("Successfully update content").
			Body("content", pageContentRes),
	)
}

// Delete deletes a PageContent by space id and page id
func (controller PageContentController) Delete(w http.ResponseWriter, r *http.Request) {
	input, err := combh.FormDataDecodeValidate[pc_reqs.Action](r.Form)
	errorh.Panic(err)

	err = controller.policy.Delete(r.Context(), input)
	errorh.Panic(err)

	// Delete page content by id
	err = controller.service.Delete(r.Context(), input)
	errorh.Panic(err)

	rwh.WriteResponse(w, responses.NewResponse().
		Code(http.StatusOK).
		Success("Successfully delete content"),
	)
}
