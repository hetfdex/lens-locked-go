package controller

const galleryIdKey = "gallery_id"

const homeRoute = "/"
const homeFilename = "view/home.gohtml"

const registerUserRoute = "/register"
const registerUserFilename = "view/user_register.gohtml"

const loginUserRoute = "/login"
const loginUserFilename = "view/user_login.gohtml"

const logoutUserRoute = "/logout"

const indexGalleryRoute = "/gallery"
const indexGalleryFilename = "view/gallery_index.gohtml"

const createGalleryRoute = "/gallery/create"
const createGalleryFilename = "view/gallery_create.gohtml"

const editGalleryRoute = "/gallery/{%s}/edit"
const editGalleryFilename = "view/gallery_edit.gohtml"

const uploadGalleryRoute = "/gallery/{%s}/upload"
const uploadGalleryFilename = "view/gallery_upload.gohtml"

const deleteGalleryRoute = "/gallery/{%s}/delete"

const galleryRoute = "/gallery/{%s}"
const galleryFilename = "view/gallery.gohtml"

const dropboxConnectRoute = "/oauth/dropbox/connect"

const dropboxCallbackRoute = "/oauth/dropbox/callback"

const queryRoute = "/oauth/dropbox/query"
