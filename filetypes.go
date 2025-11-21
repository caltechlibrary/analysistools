package analysistools

import (
	"fmt"
    "io/fs"
    "os"
	"path/filepath"
	"strings"
)

var extensionToMIME = map[string]string{
    // Text and Code
    ".txt":      "text/plain",
    ".csv":      "text/csv",
    ".html":     "text/html",
    ".htm":      "text/html",
    ".css":      "text/css",
    ".js":       "application/javascript",
    ".json":     "application/json",
    ".xml":      "application/xml",
    ".go":       "text/x-go",
    ".py":       "text/x-python",
    ".java":     "text/x-java-source",
    ".c":        "text/x-c",
    ".cpp":      "text/x-c",
    ".h":        "text/x-c",
    ".hpp":      "text/x-c",
    ".sh":       "text/x-sh",
    ".md":       "text/markdown",
    ".rtf":      "application/rtf",

    // Archives and Compressed
    ".zip":      "application/zip",
    ".tar":      "application/x-tar",
    ".gz":       "application/gzip",
    ".rar":      "application/x-rar-compressed",
    ".7z":       "application/x-7z-compressed",

    // Images
    ".png":      "image/png",
    ".jpg":      "image/jpeg",
    ".jpeg":     "image/jpeg",
    ".gif":      "image/gif",
    ".bmp":      "image/bmp",
    ".svg":      "image/svg+xml",
    ".tiff":     "image/tiff",
    ".webp":     "image/webp",
    ".ico":      "image/x-icon",

    // Documents
    ".pdf":      "application/pdf",
    ".doc":      "application/msword",
    ".docx":     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ".xls":      "application/vnd.ms-excel",
    ".xlsx":     "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ".ppt":      "application/vnd.ms-powerpoint",
    ".pptx":     "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    ".odt":      "application/vnd.oasis.opendocument.text",
    ".ods":      "application/vnd.oasis.opendocument.spreadsheet",
    ".odp":      "application/vnd.oasis.opendocument.presentation",

    // Microsoft Office (Legacy and Open XML)
    ".dot":      "application/msword",
    ".dotx":     "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
    ".xlt":      "application/vnd.ms-excel",
    ".xltx":     "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
    ".pot":      "application/vnd.ms-powerpoint",
    ".potx":     "application/vnd.openxmlformats-officedocument.presentationml.template",
    ".pps":      "application/vnd.ms-powerpoint",
    ".ppsx":     "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
    ".pub":      "application/x-mspublisher",
    ".vsdx":     "application/vnd.ms-visio.drawing.main+xml",
    ".vsd":      "application/vnd.visio",
    ".msg":      "application/vnd.ms-outlook",

    // macOS and Apple
    ".pages":    "application/x-iwork-pages-sffpages",
    ".numbers":  "application/x-iwork-numbers-sffnumbers",
    ".key":      "application/x-iwork-keynote-sffkey",
    ".dmg":      "application/x-apple-diskimage",
    ".app":      "application/x-executable",
    ".pkg":      "application/x-newton-compatible-pkg",
    ".plist":    "application/x-plist",

    // Audio
    ".mp3":      "audio/mpeg",
    ".wav":      "audio/wav",
    ".aac":      "audio/aac",
    ".m4a":      "audio/mp4",
    ".ogg":      "audio/ogg",
    ".flac":     "audio/flac",

    // Video
    ".mp4":      "video/mp4",
    ".mov":      "video/quicktime",
    ".avi":      "video/x-msvideo",
    ".wmv":      "video/x-ms-wmv",
    ".mkv":      "video/x-matroska",
    ".webm":     "video/webm",

    // Fonts
    ".woff":     "font/woff",
    ".woff2":    "font/woff2",
    ".ttf":      "font/ttf",
    ".otf":      "font/otf",

    // Executables and Binaries
    ".exe":      "application/x-msdownload",
    ".dll":      "application/x-msdownload",
    ".msi":      "application/x-msi",
    ".bat":      "application/x-msdos-program",

    // Data and Config
    ".ini":      "text/plain",
    ".conf":     "text/plain",
    ".log":      "text/plain",
    ".sql":      "application/sql",
    ".db":       "application/x-sqlite3",
    ".sqlite":   "application/x-sqlite3",

    // Web and Network
    ".wasm":     "application/wasm",
    ".jsonld":   "application/ld+json",
    ".rss":      "application/rss+xml",
    ".atom":     "application/atom+xml",

	// Add more extensions and MIME types as needed
}

func FileTypes(initialDir string, excludeList []string) (map[string]string, error) {
    var lastErr error

    fileTypes := make(map[string]string)

	err := filepath.WalkDir(initialDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %+v: %s\n", d.Type(), err)
            lastErr = err
            return nil
		}

		// Skip if it's a directory and in the exclude list
		if d.IsDir() {
			for _, exclude := range excludeList {
				if strings.Contains(path, exclude) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Skip if the file is in the exclude list
		for _, exclude := range excludeList {
			if strings.Contains(path, exclude) {
				return nil
			}
		}

		// Get file extension
		ext := strings.ToLower(filepath.Ext(path))

		// Get MIME type from the map, default to "application/octet-stream"
		mimeType, ok := extensionToMIME[ext]
		if !ok {
			mimeType = "application/octet-stream"
		}

		fileTypes[path] = mimeType
		return nil
	})
    if lastErr != nil {
        if err != nil {
        	return fileTypes, fmt.Errorf("%s\n%s\n", lastErr, err)
        }
    	return fileTypes, lastErr
    }
    if err != nil {
    	return fileTypes, err
    }
    return fileTypes, nil
}
