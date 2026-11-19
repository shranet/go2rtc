package go2rtc

import (
	http2 "net/http"

	"github.com/AlexxIT/go2rtc/internal/api"
	"github.com/AlexxIT/go2rtc/internal/api/ws"
	"github.com/AlexxIT/go2rtc/internal/app"
	"github.com/AlexxIT/go2rtc/internal/mjpeg"
	"github.com/AlexxIT/go2rtc/internal/mp4"
	"github.com/AlexxIT/go2rtc/internal/rtsp"
	"github.com/AlexxIT/go2rtc/internal/streams"
	"github.com/AlexxIT/go2rtc/pkg/tcp"
)

// InitSmartBozor — AGENT uchun minimal go2rtc (2026-11-19, xavfsizlik auditi): faqat kameradan RTSP olish va
// /api/ws orqali MSE/MJPEG ko'rsatish. exec, ffmpeg, echo, expr, webrtc, hls va boshqa manbalar/serverlar
// ulanmaydi: ular buyruq bajarish yoki ichki tarmoqqa so'rov yuborish yo'lini ochadi. Oqimlarni faqat AGENT
// yaratadi (streams.New/Patch, rtsp:// manba); /api/ws faqat mavjud oqimni ko'rsatadi (streams.GetOrPatch).
func InitSmartBozor() {
	app.Version = "1.9.11"

	app.Init() // config va loglar

	api.Init() // API (HandleFunc no-op; AGENT o'z handlerlarini Handle bilan qo'yadi)
	ws.Init()  // /api/ws

	streams.Init()
	streams.AllowSchemes("rtsp", "rtsps")
	tcp.DisableBasic = true // kameralarga parol faqat Digest bilan

	rtsp.Init()  // rtsp manba (server tinglash AGENT konfigida o'chiq)
	mp4.Init()   // MSE (ws)
	mjpeg.Init() // MJPEG (ws)
}

func Handle(pattern string, handler http2.Handler) {
	api.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler http2.HandlerFunc) {
	api.HandleFunc2(pattern, handler)
}

func NewStream(name string, sources ...string) *streams.Stream {
	return streams.New(name, sources...)
}

func UpdateStream(name string, source string) {
	streams.Patch(name, source)
}

func RemoveStream(name string) {
	streams.Delete(name)
}
