package setup

import "mime"

func MimeTypes() {
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".ico", "image/x-icon")
	mime.AddExtensionType(".gif", "image/gif")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".htm", "text/html")
	mime.AddExtensionType(".html", "text/html")
	mime.AddExtensionType(".png", "image/png")
}
