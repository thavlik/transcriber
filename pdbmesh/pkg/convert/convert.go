package convert

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
	"github.com/google/uuid"

	"github.com/pkg/errors"
)

type Model struct {
	f    *os.File
	path string
}

func (m *Model) Reader() io.Reader {
	return m.f
}

func (m *Model) Dispose() error {
	defer os.Remove(m.path)
	if err := m.f.Close(); err != nil {
		return errors.Wrap(err, "failed to close file")
	}
	return nil
}

func Convert(
	ctx context.Context,
	r io.Reader,
) (*Model, error) {
	models, err := pdb.NewReader(r).ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PDB file")
	} else if len(models) == 0 {
		return nil, errors.New("no models found in PDB file")
	}
	model := models[0]
	mesh := ribbon.ModelMesh(model)
	id := uuid.New().String()
	outPath := filepath.Join("/tmp", fmt.Sprintf("%s.stl", id))
	if err := mesh.SaveSTL(outPath); err != nil {
		return nil, errors.Wrap(err, "failed to save STL file")
	}
	f, err := os.Open(outPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open STL file")
	}
	return &Model{
		f:    f,
		path: outPath,
	}, nil
}
