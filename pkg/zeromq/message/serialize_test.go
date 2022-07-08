package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/streambot/pkg/zeromq/message"
)

func TestSerialization(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		var valueBefore = "hello"
		var valueAfter string
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("int", func(t *testing.T) {
		var valueBefore = 5
		var valueAfter int
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("float64", func(t *testing.T) {
		var valueBefore = 5.5
		var valueAfter float64
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("[]string", func(t *testing.T) {
		var valueBefore = []string{"a", "b"}
		var valueAfter []string
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})

	t.Run("map", func(t *testing.T) {
		var valueBefore = map[string]int{"foo": 2}
		var valueAfter map[string]int
		message.Unserialize(message.Serialize(valueBefore), &valueAfter)
		assert.Equal(t, valueBefore, valueAfter)
	})
}

func BenchmarkSerialize(b *testing.B) {
	b.ReportAllocs()

	serialize := message.Serialize
	unserialize := message.Unserialize

	b.Run("short string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			valueBefore := "Lorem ipsum"
			var valueAfter string
			unserialize(serialize(valueBefore), &valueAfter)
		}
	})

	b.Run("long string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			valueBefore := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam sed ipsum sed diam scelerisque fermentum. Fusce sed iaculis sapien, at aliquam diam. Sed id urna enim. Morbi non tempor erat. Pellentesque at lacinia libero. Proin molestie, nibh in maximus placerat, mauris ante rutrum nulla, a scelerisque tellus mi a mi. Nunc sed est mauris."
			var valueAfter string
			unserialize(serialize(valueBefore), &valueAfter)
		}
	})

	b.Run("[]string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			valueBefore := []string{"alpha", "beta", "gamma", "delta"}
			var valueAfter []string
			unserialize(serialize(valueBefore), &valueAfter)
		}
	})

	b.Run("generic server", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			valueBefore := qserver.GenericServer{
				Version: qversion.New("MVDSV 0.35-dev"),
				Address: ":8001",
				Clients: []qclient.Client{
					{
						Name:   qstring.New("NL"),
						Team:   qstring.New("red"),
						Skin:   "",
						Colors: [2]uint8{13, 13},
						Frags:  2,
						Ping:   38,
						Time:   4,
					},
					{
						Name:   qstring.New("[ServeMe]"),
						Team:   qstring.New("lqwc"),
						Skin:   "",
						Colors: [2]uint8{12, 11},
						Frags:  -9999,
						Ping:   -666,
						Time:   16,
					},
				},
				Settings: map[string]string{
					"*version":        "MVDSV 0.35-dev",
					"hostname":        "troopers.fi:28501\u0087",
					"hostname_parsed": "troopers.fi:28501",
					"maxfps":          "77",
					"pm_ktjump":       "1",
				},
				Geo: geo.Info{},
				ExtraInfo: struct {
					QtvStream qtvstream.QtvStream `json:"qtv_stream"`
				}{
					QtvStream: qtvstream.QtvStream{
						SpectatorNames: make([]string, 0),
					},
				},
			}

			var valueAfter qserver.GenericServer
			unserialize(serialize(valueBefore), &valueAfter)
		}
	})
}

func TestSerializedValue_ToString(t *testing.T) {
	assert.Equal(t, "abc", message.Serialize("abc").ToString())
}

func TestSerializedValue_ToInt(t *testing.T) {
	assert.Equal(t, 123, message.Serialize(123).ToInt())
}

func TestSerializedValue_To(t *testing.T) {
	valueBefore := []string{"a", "b", "c"}
	var valueAfter []string
	message.Serialize(valueBefore).To(&valueAfter)
	assert.Equal(t, valueBefore, valueAfter)
}
