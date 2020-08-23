package entity_test

import (
	"go/ast"
	"reflect"
	"testing"

	"github.com/Nao-Mk2/go-class-diagram-generator/internal/entity"
)

func TestNewPackage(t *testing.T) {
	type args struct {
		name    string
		imports []*ast.ImportSpec
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Package
		wantErr bool
	}{
		{
			name: "return a argument error",
			args: args{
				name: "",
			},
			wantErr: true,
		},
		{
			name: "return a Package without Imports",
			args: args{
				name: "main",
			},
			want: &entity.Package{
				Name:    "main",
				Imports: make([]entity.Import, 0, 0),
			},
		},
		{
			name: "return a Package with Imports",
			args: args{
				name: "main",
				imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{
							Value: "fmt",
						},
					},
					{
						Path: &ast.BasicLit{
							Value: "github.com/Nao-Mk2/go-class-diagram-generator/internal/parser",
						},
					},
				},
			},
			want: &entity.Package{
				Name: "main",
				Imports: []entity.Import{
					{
						Path:              "fmt",
						IsStandardLibrary: true,
					},
					{
						Path:              "github.com/Nao-Mk2/go-class-diagram-generator/internal/parser",
						IsStandardLibrary: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewPackage(tt.args.name, tt.args.imports)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}
