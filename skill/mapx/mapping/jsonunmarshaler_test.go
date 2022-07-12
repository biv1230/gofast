package mapping

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalBytes(t *testing.T) {
	var c struct {
		Name string
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalBytesOptional(t *testing.T) {
	var c struct {
		Name string
		Age  int `cnf:",NA"`
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalBytesOptionalDefault(t *testing.T) {
	var c struct {
		Name string
		Age  int `cnf:",NA,def=1"`
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalBytesDefaultOptional(t *testing.T) {
	var c struct {
		Name string
		Age  int `cnf:",def=1,NA"`
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalBytesDefault(t *testing.T) {
	var c struct {
		Name string `cnf:",def=liao"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalBytesBool(t *testing.T) {
	var c struct {
		Great bool
	}
	content := []byte(`{"Great": true}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.True(t, c.Great)
}

func TestUnmarshalBytesInt(t *testing.T) {
	var c struct {
		Age int
	}
	content := []byte(`{"Age": 1}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalBytesUint(t *testing.T) {
	var c struct {
		Age uint
	}
	content := []byte(`{"Age": 1}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, uint(1), c.Age)
}

func TestUnmarshalBytesFloat(t *testing.T) {
	var c struct {
		Age float32
	}
	content := []byte(`{"Age": 1.5}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, float32(1.5), c.Age)
}

func TestUnmarshalBytesMustInOptional(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `cnf:",NA"`
		} `cnf:",NA"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMustInOptionalMissedPart(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `cnf:",NA"`
		} `cnf:",NA"`
	}
	content := []byte(`{"Inner": {"There": "sure"}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMustInOptionalOnlyOptionalFilled(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `cnf:",NA"`
		} `cnf:",NA"`
	}
	content := []byte(`{"Inner": {"Optional": "sure"}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesNil(t *testing.T) {
	var c struct {
		Int int64 `cnf:"int,NA"`
	}
	content := []byte(`{"int":null}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, int64(0), c.Int)
}

func TestUnmarshalBytesNilSlice(t *testing.T) {
	var c struct {
		Ints []int64 `cnf:"ints"`
	}
	content := []byte(`{"ints":[null]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 0, len(c.Ints))
}

func TestUnmarshalBytesPartial(t *testing.T) {
	var c struct {
		Name string
		Age  float32
	}
	content := []byte(`{"Age": 1.5}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesStruct(t *testing.T) {
	var c struct {
		Inner struct {
			Name string
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalBytesStructOptional(t *testing.T) {
	var c struct {
		Inner struct {
			Name string
			Age  int `cnf:",NA"`
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalBytesStructPtr(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalBytesStructPtrOptional(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
			Age  int `cnf:",NA"`
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesStructPtrDefault(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
			Age  int `cnf:",def=4"`
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
	assert.Equal(t, 4, c.Inner.Age)
}

func TestUnmarshalBytesSliceString(t *testing.T) {
	var c struct {
		Names []string
	}
	content := []byte(`{"Names": ["liao", "chaoxin"]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(c.Names, want) {
		t.Fatalf("want %q, got %q", c.Names, want)
	}
}

func TestUnmarshalBytesSliceAttrString(t *testing.T) {
	var c struct {
		Names []string
		Age   []int `cnf:",NA"`
	}
	content := []byte(`{"Names": ["liao", "chaoxin"]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(c.Names, want) {
		t.Fatalf("want %q, got %q", c.Names, want)
	}
}

func TestUnmarshalBytesSliceStruct(t *testing.T) {
	var c struct {
		People []struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []struct {
		Name string
		Age  int
	}{
		{"liao", 1},
		{"chaoxin", 2},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %q, got %q", c.People, want)
	}
}

func TestUnmarshalBytesSliceStructOptional(t *testing.T) {
	var c struct {
		People []struct {
			Name   string
			Age    int
			Emails []string `cnf:",NA"`
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []struct {
		Name   string
		Age    int
		Emails []string `cnf:",NA"`
	}{
		{"liao", 1, nil},
		{"chaoxin", 2, nil},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %q, got %q", c.People, want)
	}
}

func TestUnmarshalBytesSliceStructPtr(t *testing.T) {
	var c struct {
		People []*struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []*struct {
		Name string
		Age  int
	}{
		{"liao", 1},
		{"chaoxin", 2},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %v, got %v", c.People, want)
	}
}

func TestUnmarshalBytesSliceStructPtrOptional(t *testing.T) {
	var c struct {
		People []*struct {
			Name   string
			Age    int
			Emails []string `cnf:",NA"`
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []*struct {
		Name   string
		Age    int
		Emails []string `cnf:",NA"`
	}{
		{"liao", 1, nil},
		{"chaoxin", 2, nil},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %v, got %v", c.People, want)
	}
}

func TestUnmarshalBytesSliceStructPtrPartial(t *testing.T) {
	var c struct {
		People []*struct {
			Name  string
			Age   int
			Email string
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesSliceStructPtrDefault(t *testing.T) {
	var c struct {
		People []*struct {
			Name  string
			Age   int
			Email string `cnf:",def=chaoxin@liao.com"`
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []*struct {
		Name  string
		Age   int
		Email string
	}{
		{"liao", 1, "chaoxin@liao.com"},
		{"chaoxin", 2, "chaoxin@liao.com"},
	}

	for i := range c.People {
		actual := c.People[i]
		expect := want[i]
		assert.Equal(t, expect.Age, actual.Age)
		assert.Equal(t, expect.Email, actual.Email)
		assert.Equal(t, expect.Name, actual.Name)
	}
}

func TestUnmarshalBytesSliceStringPartial(t *testing.T) {
	var c struct {
		Names []string
		Age   int
	}
	content := []byte(`{"Age": 1}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesSliceStructPartial(t *testing.T) {
	var c struct {
		Group  string
		People []struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"Group": "chaoxin"}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesInnerAnonymousPartial(t *testing.T) {
	type (
		Deep struct {
			A string
			B string `cnf:",NA"`
		}
		Inner struct {
			Deep
			InnerV string `cnf:",NA"`
		}
	)

	var c struct {
		Value Inner `cnf:",NA"`
	}
	content := []byte(`{"Value": {"InnerV": "chaoxin"}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesStructPartial(t *testing.T) {
	var c struct {
		Group  string
		Person struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"Group": "chaoxin"}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesEmptyMap(t *testing.T) {
	var c struct {
		Persons map[string]int `cnf:",NA"`
	}
	content := []byte(`{"Persons": {}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Empty(t, c.Persons)
}

func TestUnmarshalBytesMap(t *testing.T) {
	var c struct {
		Persons map[string]int
	}
	content := []byte(`{"Persons": {"first": 1, "second": 2}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 2, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"])
	assert.Equal(t, 2, c.Persons["second"])
}

func TestUnmarshalBytesMapStruct(t *testing.T) {
	var c struct {
		Persons map[string]struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": {"Id": 1, "name": "kevin"}}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].Id)
	assert.Equal(t, "kevin", c.Persons["first"].Name)
}

func TestUnmarshalBytesMapStructPtr(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": {"Id": 1, "name": "kevin"}}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].Id)
	assert.Equal(t, "kevin", c.Persons["first"].Name)
}

func TestUnmarshalBytesMapStructMissingPartial(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			Id   int
			Name string
		}
	}
	content := []byte(`{"Persons": {"first": {"Id": 1}}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMapStructOptional(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": {"Id": 1}}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].Id)
}

func TestUnmarshalBytesMapEmptyStructSlice(t *testing.T) {
	var c struct {
		Persons map[string][]struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": []}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Empty(t, c.Persons["first"])
}

func TestUnmarshalBytesMapStructSlice(t *testing.T) {
	var c struct {
		Persons map[string][]struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": [{"Id": 1, "name": "kevin"}]}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].Id)
	assert.Equal(t, "kevin", c.Persons["first"][0].Name)
}

func TestUnmarshalBytesMapEmptyStructPtrSlice(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": []}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Empty(t, c.Persons["first"])
}

func TestUnmarshalBytesMapStructPtrSlice(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": [{"Id": 1, "name": "kevin"}]}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].Id)
	assert.Equal(t, "kevin", c.Persons["first"][0].Name)
}

func TestUnmarshalBytesMapStructPtrSliceMissingPartial(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			Id   int
			Name string
		}
	}
	content := []byte(`{"Persons": {"first": [{"Id": 1}]}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMapStructPtrSliceOptional(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			Id   int
			Name string `cnf:"name,NA"`
		}
	}
	content := []byte(`{"Persons": {"first": [{"Id": 1}]}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].Id)
}

func TestUnmarshalStructOptional(t *testing.T) {
	var c struct {
		Name string
		Etcd struct {
			Hosts []string
			Key   string
		} `cnf:",NA"`
	}
	content := []byte(`{"Name": "kevin"}`)

	err := UnmarshalJsonBytes(content, &c)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", c.Name)
}

func TestUnmarshalStructLowerCase(t *testing.T) {
	var c struct {
		Name string
		Etcd struct {
			Key string
		} `cnf:"etcd"`
	}
	content := []byte(`{"Name": "kevin", "etcd": {"Key": "the key"}}`)

	err := UnmarshalJsonBytes(content, &c)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", c.Name)
	assert.Equal(t, "the key", c.Etcd.Key)
}

func TestUnmarshalWithStructAllOptionalWithEmpty(t *testing.T) {
	var c struct {
		Inner struct {
			Optional string `cnf:",NA"`
		}
		Else string
	}
	content := []byte(`{"Else": "sure", "Inner": {}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructAllOptionalPtr(t *testing.T) {
	var c struct {
		Inner *struct {
			Optional string `cnf:",NA"`
		}
		Else string
	}
	content := []byte(`{"Else": "sure", "Inner": {}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		In   Inner `cnf:",NA"`
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "", c.In.Must)
}

func TestUnmarshalWithStructPtrOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		In   *Inner `cnf:",NA"`
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Nil(t, c.In)
}

func TestUnmarshalWithStructAllOptionalAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `cnf:",NA"`
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructAllOptionalAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `cnf:",NA"`
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructAllOptionalProvoidedAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `cnf:",NA"`
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Optional": "optional"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "optional", c.Optional)
}

func TestUnmarshalWithStructAllOptionalProvoidedAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `cnf:",NA"`
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Optional": "optional"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "optional", c.Optional)
}

func TestUnmarshalWithStructAnonymous(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Must": "must"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "must", c.Must)
}

func TestUnmarshalWithStructAnonymousPtr(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Must": "must"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "must", c.Must)
}

func TestUnmarshalWithStructAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		Inner `cnf:",NA"`
		Else  string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "", c.Must)
}

func TestUnmarshalWithStructPtrAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		*Inner `cnf:",NA"`
		Else   string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Nil(t, c.Inner)
}

func TestUnmarshalWithZeroValues(t *testing.T) {
	type inner struct {
		False  bool   `cnf:"no"`
		Int    int    `cnf:"int"`
		String string `cnf:"string"`
	}
	content := []byte(`{"no": false, "int": 0, "string": ""}`)
	reader := bytes.NewReader(content)

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalJsonReader(reader, &in))
	ast.False(in.False)
	ast.Equal(0, in.Int)
	ast.Equal("", in.String)
}

func TestUnmarshalBytesError(t *testing.T) {
	payload := `[{"abcd": "cdef"}]`
	var v struct {
		Any string
	}

	err := UnmarshalJsonBytes([]byte(payload), &v)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), payload))
}

func TestUnmarshalReaderError(t *testing.T) {
	payload := `[{"abcd": "cdef"}]`
	reader := strings.NewReader(payload)
	var v struct {
		Any string
	}

	err := UnmarshalJsonReader(reader, &v)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), payload))
}