package streams

import (
	"net/url"
	"testing"

	"github.com/AlexxIT/go2rtc/pkg/core"
)

// smartbozor: /api/ws?src=exec:... yangi oqim (va buyruq) yarata olmaydi; manba faqat rtsp.
func TestSmartBozorNoDynamicSources(t *testing.T) {
	HandleFunc("exec", func(string) (core.Producer, error) { t.Fatal("exec chaqirilmasligi kerak"); return nil, nil })
	HandleFunc("rtsp", func(string) (core.Producer, error) { return nil, nil })
	AllowSchemes("rtsp", "rtsps")
	defer func() {
		allowedSchemes = nil
		streamsMu.Lock()
		for k := range streams {
			delete(streams, k)
		}
		streamsMu.Unlock()
	}()

	for _, src := range []string{"exec:touch /tmp/x", "exec:id", "ffmpeg:rtsp://a/b", "rtsp://10.0.0.1/x"} {
		if s := GetOrPatch(url.Values{"src": {src}}); s != nil {
			t.Fatalf("dinamik oqim yaratildi: %s", src)
		}
		if s := GetOrPatch(url.Values{"src": {src}, "name": {"n"}}); s != nil {
			t.Fatalf("name bilan dinamik oqim yaratildi: %s", src)
		}
	}
	if New("x", "exec:id") != nil || Patch("y", "exec:id") != nil {
		t.Fatal("exec manbasi qabul qilinmasligi kerak")
	}
	if New("cam", "rtsp://u:p@10.0.0.1:554/Streaming/Channels/101") == nil {
		t.Fatal("rtsp manba ishlashi kerak")
	}
	if GetOrPatch(url.Values{"src": {"cam"}}) == nil {
		t.Fatal("mavjud oqim topilishi kerak")
	}
}

func TestRedactURL(t *testing.T) {
	if got := redactURL("rtsp://admin:S3cr%23t@192.168.1.4:554/Streaming/channels/101"); got !=
		"rtsp://%2A%2A%2A@192.168.1.4:554/Streaming/channels/101" && got != "rtsp://***@192.168.1.4:554/Streaming/channels/101" {
		t.Fatalf("redact: %s", got)
	}
	if got := redactURL("rtsp://192.168.1.4/x"); got != "rtsp://192.168.1.4/x" {
		t.Fatalf("userinfo siz: %s", got)
	}
}
