package treelite_test

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/saikathowladarzomato/go-treelite"
)

func Example() {
	data, nRow, nCol := loadData()

	dMatrix, err := treelite.CreateFromMat(data, nRow, nCol, float32(math.NaN()))
	if err != nil {
		log.Fatal(err)
	}
	defer dMatrix.Close()

	model, err := treelite.LoadXGBoostModel("testdata/xgboost.model")
	if err != nil {
		log.Fatal(err)
	}
	defer model.Close()

	annotator, err := treelite.NewAnnotator(model, dMatrix, 1, true)
	if err != nil {
		log.Fatal(err)
	}
	defer annotator.Close()

	err = annotator.Save("testdata/go-example-annotation.json")
	if err != nil {
		log.Fatal(err)
	}

	compiler, err := treelite.NewCompiler(
		"ast_native",
		&treelite.CompilerParam{
			AnnotationPath: "testdata/go-example-annotation.json",
			Quantize:       true,
			ParallelComp:   runtime.NumCPU(),
			Verbose:        true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer compiler.Close()

	err = compiler.ExportSharedLib(
		model,
		"testdata/go_example_compiled_model",
		"gcc",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	predictor, err := treelite.NewPredictor(
		fmt.Sprintf("testdata/go_example_compiled_model.%s", treelite.GetSharedLibExtension()),
		runtime.NumCPU(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer predictor.Close()

	scores, err := predictor.PredictBatch(dMatrix, true, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", scores)
}

func loadData() ([]float32, int, int) {
	nCol := 30
	var nRow int
	feature := make([]float32, 0)
	featureFile, err := os.Open("testdata/feature.csv")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(featureFile)
	for scanner.Scan() {
		nRow++
		featureValues := strings.Split(scanner.Text(), ",")
		for _, valueString := range featureValues {
			value, err := strconv.ParseFloat(valueString, 32)
			if err != nil {
				log.Fatal(err)
			}
			feature = append(feature, float32(value))
		}
	}

	return feature, nRow, nCol
}

func TestEndToEnd_Example(t *testing.T) {
	Example()
}
