package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sipropeda-backend/internal/domain/auth"
	"sipropeda-backend/shared/failure"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type MenuHandler struct {
	MenuService auth.MenuService
}

func ProvideMenuHandler(MenuService auth.MenuService) MenuHandler {
	return MenuHandler{
		MenuService: MenuService,
	}
}

func (h *MenuHandler) Router(r chi.Router) {
	r.Route("/menu-role", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected) 
			protected.Get("/", h.ResolveMenuByRoleID)
			protected.Get("/trx", h.ResolveMenuByRoleIDTrx)
			protected.Post("/bulk", h.CreateBulkMenuRole)
			protected.Put("/update-permission", h.UpdateMenuPermission)
		})
	})
	r.Route("/menu", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Get("/", h.ResolveAll)
			protected.Get("/all", h.GetAllMenu)
			protected.Get("/{id}", h.ResolveMenuByID)
			protected.Post("/", h.CreateMenu)
			protected.Put("/{id}", h.UpdateMenu)
			protected.Delete("/{id}", h.DeleteMenu)
		})
	})
}

// GetAllMenu mengambil semua data Menu tanpa pagination
// @Summary Ambil Semua Menu
// @Tags Menu
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/menu/all [get]
func (h *MenuHandler) GetAllMenu(w http.ResponseWriter, r *http.Request) {
	status, err := h.MenuService.GetAllMenu()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, status)
}

// ResolveAll mengambil semua data Menu dengan fitur pencarian dan pagination
// @Summary Ambil Semua Menu (Pagination)
// @Tags Menu
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Keyword Pencarian"
// @Param pageSize query int false "Ukuran Halaman"
// @Param pageNumber query int false "Nomor Halaman"
// @Param sortBy query string false "Sort by field"
// @Param sortType query string false "Sort type (ASC/DESC)"
// @Success 200 {object} response.Base
// @Router /v1/menu [get]
func (h *MenuHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "DESC"
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	req := model.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
	}

	status, err := h.MenuService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// CreateMenu menambah data Menu baru
// @Summary Tambah Menu Baru
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestMenuFormat true "Data Menu"
// @Success 201 {object} response.Base
// @Router /v1/menu [post]
func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestMenuFormat
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	// Buat UUID otomatis untuk system admin/user (Jika tidak pakai GetClaimsValue)
	userID, _ := uuid.NewV4() 
	reqFormat.UserID = userID

	newMenu, err := h.MenuService.CreateMenu(reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenu)
}

// ResolveMenuByID mengambil detail data Menu berdasarkan ID
// @Summary Ambil detail Menu by ID
// @Tags Menu
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Menu"
// @Success 200 {object} response.Base
// @Router /v1/menu/{id} [get]
func (h *MenuHandler) ResolveMenuByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	menu, err := h.MenuService.ResolveMenuByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, menu)
}

// UpdateMenu memperbarui data Menu berdasarkan ID
// @Summary Perbarui data Menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Menu"
// @Param body body auth.RequestMenuFormat true "Data Menu yang diperbarui"
// @Success 200 {object} response.Base
// @Router /v1/menu/{id} [put]
func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
        return
	}

	var newMenu auth.RequestMenuFormat
	err = json.NewDecoder(r.Body).Decode(&newMenu)
    if err != nil {
		response.WithError(w, failure.BadRequest(err))
        return
	}

	userID, _ := uuid.NewV4()
	newMenu.UserID = userID

	menu, err := h.MenuService.UpdateMenu(id, newMenu)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// DeleteMenu menghapus data Menu berdasarkan ID
// @Summary Hapus Menu by ID
// @Tags Menu
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Menu"
// @Success 200 {object} response.Base
// @Router /v1/menu/{id} [delete]
func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newID, err := uuid.FromString(id)
    if err != nil {
		response.WithError(w, failure.BadRequest(err))
        return
	}

	err = h.MenuService.DeleteMenuByID(newID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Data Berhasil di Hapus",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// ResolveMenuByRoleID mengambil data Menu berdasarkan Role ID
// @Summary Ambil Menu by Role ID
// @Tags Menu Role
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param roleId query string true "ID Role"
// @Param commodityId query string false "ID Commodity"
// @Success 200 {object} response.Base
// @Router /v1/menu-role [get]
func (h *MenuHandler) ResolveMenuByRoleID(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	commodityID := r.URL.Query().Get("commodityId")
	req := auth.MenuRequest{
		RoleId:      roleID,
		CommodityId: commodityID,
	}

	menu, err := h.MenuService.ResolveMenuByRoleID(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// ResolveMenuByRoleIDTrx mengambil data Menu Transaksi berdasarkan Role ID
// @Summary Ambil Menu Transaksi by Role ID
// @Tags Menu Role
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param roleId query string true "ID Role"
// @Param commodityId query string false "ID Commodity"
// @Success 200 {object} response.Base
// @Router /v1/menu-role/trx [get]
func (h *MenuHandler) ResolveMenuByRoleIDTrx(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	commodityID := r.URL.Query().Get("commodityId")
	req := auth.MenuRequest{
		RoleId:      roleID,
		CommodityId: commodityID,
	}

	menu, err := h.MenuService.ResolveMenuByRoleIDTrx(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// UpdateMenuPermission memperbarui izin/permission suatu menu
// @Summary Perbarui Menu Permission
// @Tags Menu Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestMenuPermissionFormat true "Data Menu Permission"
// @Success 201 {object} response.Base
// @Router /v1/menu-role/update-permission [put]
func (h *MenuHandler) UpdateMenuPermission(w http.ResponseWriter, r *http.Request) {
	var menu auth.RequestMenuPermissionFormat
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = h.MenuService.UpdateMenuPermission(menu)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, "success")
}

// CreateBulkMenuRole menambahkan Menu Role secara massal
// @Summary Tambah Bulk Menu Role
// @Tags Menu Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestBulkMenuRole true "Data Bulk Menu Role"
// @Success 201 {object} response.Base
// @Router /v1/menu-role/bulk [post]
func (h *MenuHandler) CreateBulkMenuRole(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestBulkMenuRole
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newMenuUser, err := h.MenuService.CreateBulkMenuRole(reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenuUser)
}