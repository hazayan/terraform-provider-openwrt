package host

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestReadResponseOptionMACAddress(t *testing.T) {
	t.Run("reads UCI list values into the existing string attribute", func(t *testing.T) {
		// Given
		section := lucirpc.Options{
			"mac": lucirpc.ListString([]string{
				"12:34:56:78:90:ab",
				"ab:90:78:56:34:12",
			}),
		}

		// When
		_, got, diagnostics := readResponseOptionMACAddress(context.Background(), "openwrt", "resource", section, model{})

		// Then
		assert.Assert(t, !diagnostics.HasError())
		assert.Equal(t, got.MACAddress.ValueString(), "12:34:56:78:90:ab ab:90:78:56:34:12")
	})

	t.Run("keeps absent MAC values null", func(t *testing.T) {
		// When
		_, got, diagnostics := readResponseOptionMACAddress(context.Background(), "openwrt", "resource", lucirpc.Options{}, model{})

		// Then
		assert.Assert(t, !diagnostics.HasError())
		assert.Assert(t, got.MACAddress.IsNull())
	})
}

func TestUpsertRequestOptionMACAddress(t *testing.T) {
	t.Run("writes whitespace-separated MAC addresses as a UCI list", func(t *testing.T) {
		// Given
		resource := model{
			MACAddress: types.StringValue("12:34:56:78:90:ab ab:90:78:56:34:12"),
		}

		// When
		_, got, diagnostics := upsertRequestOptionMACAddress(context.Background(), "openwrt", lucirpc.Options{}, resource)

		// Then
		assert.Assert(t, !diagnostics.HasError())
		gotMACs, err := got["mac"].AsListString()
		assert.NilError(t, err)
		assert.DeepEqual(t, gotMACs, []string{
			"12:34:56:78:90:ab",
			"ab:90:78:56:34:12",
		})
	})
}
