package musicsnap

//
//func (h MusicsnapHandler) GetBanner(ginCtx *gin.Context, params generated.GetBannerParams) {
//	tr := global.Tracer(domain.ServiceName)
//	ctxTrace, span := tr.Start(ginCtx, h.spanName("GetBanner"))
//	defer span.End()
//
//	ctx := zapctx.WithLogger(ctxTrace, h.logger)
//
//	actor, err := NewActorFromToken(ginCtx, params.Token)
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, err)
//		return
//	}
//
//	banners, err := h.s.GetList(ctx, actor, params.ToValidDomain())
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, err)
//		return
//	}
//
//	if banners == nil {
//		err := app.NewError(http.StatusInternalServerError, "nil banners without error", "get banners returned nil lesson without error", nil)
//		h.logger.Error("nil lesson", zap.Error(err))
//		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err)
//		return
//	}
//
//	resp := models.ToBannerListResponse(banners)
//
//	ginCtx.JSON(http.StatusOK, resp)
//}
//
//func (h MusicsnapHandler) PostBanner(ginCtx *gin.Context, params generated.PostBannerParams) {
//	tr := global.Tracer(domain.ServiceName)
//	ctxTrace, span := tr.Start(ginCtx, h.spanName("PostBanner"))
//	defer span.End()
//
//	ctx := zapctx.WithLogger(ctxTrace, h.logger)
//
//	actor, err := NewActorFromToken(ginCtx, params.Token)
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, err)
//		return
//	}
//
//	var payload generated.PostBannerJSONBody
//	h.bindRequestBody(ginCtx, &payload)
//
//	bannerPayload, err := payload.ToValidDomain()
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
//		return
//	}
//
//	bannerID, err := h.s.Create(ctx, actor, bannerPayload)
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, err)
//		return
//	}
//
//	resp := models.ToBannerResponse(domain.Banner{ID: bannerID})
//
//	ginCtx.JSON(http.StatusOK, resp)
//}
//
//func (h MusicsnapHandler) DeleteBannerId(ginCtx *gin.Context, id int, params generated.DeleteBannerIdParams) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (h MusicsnapHandler) PatchBannerId(ginCtx *gin.Context, id int, params generated.PatchBannerIdParams) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (h MusicsnapHandler) GetUserBanner(ginCtx *gin.Context, params generated.GetUserBannerParams) {
//	//TODO implement me
//	panic("implement me")
//}

//
//func (h MusicsnapHandler) GetCoursesEditList(ginCtx *gin.Context, params generated.GetCoursesEditListParams) {
//	tr := global.Tracer(domain.ServiceName)
//	ctxTrace, span := tr.Start(ginCtx, "musicsnap/handler.GetCoursesEditList")
//	defer span.End()
//
//	ctx := zapctx.WithLogger(ctxTrace, h.logger)
//
//	roles, err := h.musicsnap.GetActorRoles(ctx, params.Actor.ToValidDomain())
//
//	course, err := h.musicsnap.GetAllCoursesTemplates(ctx, params.Actor.ToDomainWithRoles(roles), params.Offset, params.Limit)
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, err)
//		return
//	}
//
//	if course == nil {
//		err := app.NewError(http.StatusInternalServerError, "nil musicsnap without error", "GetCourse returned nil musicsnap without error", nil)
//		h.logger.Error("nil course", zap.Error(err))
//		http2.AbortWithBadResponse(ginCtx, h.logger, http2.MapErrorToCode(err), err)
//	}
//	resp := models.ToCourseListResponse(course)
//
//	ginCtx.JSON(http.StatusOK, resp)
//}
//
//// PostCoursesEdit CREATE course
//func (h MusicsnapHandler) PostCoursesEdit(ginCtx *gin.Context, params generated.PostCoursesEditParams) {
//	tr := global.Tracer(domain.ServiceName)
//	ctxTrace, span := tr.Start(ginCtx, "musicsnap/handler.PostCoursesEdit")
//	defer span.End()
//
//	ctx := zapctx.WithLogger(ctxTrace, h.logger)
//
//	var coursePayload generated.Course
//	h.bindRequestBody(ginCtx, &coursePayload)
//
//	roles, err := h.musicsnap.GetActorRoles(ctx, params.Actor.ToValidDomain())
//
//	course, err := h.musicsnap.CreateCourse(ctx, params.Actor.ToDomainWithRoles(roles), coursePayload.ToValidDomain())
//	if err != nil {
//		h.abortWithAutoResponse(ginCtx, err)
//		return
//	}
//
//	if course == nil {
//		err := app.NewError(http.StatusInternalServerError, "nil course without error", "GetCourse returned nil course without error", nil)
//		h.logger.Error("nil course", zap.Error(err))
//		http2.AbortWithBadResponse(ginCtx, h.logger, http2.MapErrorToCode(err), err)
//	}
//	resp := models.ToCourseResponse(*course)
//
//	ginCtx.JSON(http.StatusOK, resp)
//}
