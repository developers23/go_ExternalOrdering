package main
import(
	"fmt"
	"os"
	"github.com/myproject/bingfa/pipeline"
	"bufio"
	)
func main(){
	infile := "small.in" 
	outfile := "small.out"
	p := createPipeline(infile,512,4)
	writeToFile(p,outfile)
	//printFile(outfile)
}

func printFile(filename string){
	file , err := os.Open(filename)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file,-1)

	for v := range p{
		fmt.Println(v)
	}
}

func writeToFile(p <-chan int , filename string){
	file ,err := os.Create(filename)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer,p)
}

func createPipeline(filename string ,fileSize,chunkCount int) <-chan int{
	chunkSize := fileSize / chunkCount
	pipeline.Init()
	sortResults := []<-chan int{}
	for i := 0 ; i < chunkCount ; i++{
		file ,err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize),0)
		source := pipeline.ReaderSource(
			bufio.NewReader(file),chunkSize)
		sortResults = append(sortResults,
			pipeline.InMemSort(source))
	}
	return pipeline.MergeN(sortResults...)
}