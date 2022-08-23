package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/iancoleman/orderedmap"
)

func compress_main(fromFile, toFile string) {
	fmt.Println("Must specify file to be compressed as first command line parameter.")
	fmt.Println("Must specify new file to be compressed into as second command line parameter.")
	fmt.Println("*****************************************")
	fmt.Println("*            Usage                 ")
	fmt.Println("* compress_main infile outfile.cmp ")
	fmt.Println("*****************************************")
	fmt.Println("\nCompressing ->", fromFile, " ->", toFile)
	fmt.Println()

	hash := initFrequencyHash(fromFile)

	encodingHash := orderedmap.New() //map[string]string{}

	initBinaryTree(&hash, encodingHash)

	originalTextBytes, err := ReadInBytesFromFile(fromFile)
	if err != nil {
		log.Fatal(err)
	}
	originalText := string(originalTextBytes)

	compressedText, lengthOfCompressedTextInBytes, lengthOfCompressedTextInBits := compressText(encodingHash, originalText)

	// Marshall - initFrequencyHash returns the hash with Node as value that are nil
	hashForDecompression := initFrequencyHashWithFloat64ForValues(fromFile)

	// write this to fileInBytesInMemory
	hashMarshalled, err := json.Marshal(hashForDecompression)
	if err != nil {
		panic(err)
	}

	// write this to fileInBytesInMemory
	marshalledHashDecompressionLength := len(hashMarshalled)

	// use this for lenght of file in bytes
	byteLengthOfCompressedTextWithAdditional := uint64(8.0 + 8.0 + uint64(marshalledHashDecompressionLength) + 8.0) // add 8.0 bytes for this size byteLengthOfCompressedText and add 8.0 for length of marshalledHashDecompressionLength (8) plus lenght of hashMarshalled and plus length of compressed text in bits (8)

	// This is actually the contents of the file (write it to the file).
	fileInBytesInMemory := make([]byte, byteLengthOfCompressedTextWithAdditional)

	marshalledHashDecompressionLengthMarshalled := getBytesForInt(marshalledHashDecompressionLength)
	// hashMarshalled
	lengthOfCompressedTextInBytesMarshalled := getBytesForInt(int(lengthOfCompressedTextInBytes))
	lengthOfCompressedTextInBitsMarshalled := getBytesForInt(int(lengthOfCompressedTextInBits))
	// compressedTextAsByteRay

	count := 0
	for count < 8 {
		fileInBytesInMemory[count] = marshalledHashDecompressionLengthMarshalled[count]
		count++
	}

	for count < 8+len(hashMarshalled) {
		fileInBytesInMemory[count] = hashMarshalled[count-8]
		count++
	}

	for count < 8+len(hashMarshalled)+8 {
		fileInBytesInMemory[count] = lengthOfCompressedTextInBytesMarshalled[count-8-len(hashMarshalled)]
		count++
	}

	for count < 8+len(hashMarshalled)+8+8 {
		fileInBytesInMemory[count] = lengthOfCompressedTextInBitsMarshalled[count-8-8-len(hashMarshalled)]
		count++
	}

	compressedTextAsByteRay := InitNewByteset(append(fileInBytesInMemory, compressedText...))

	// Open a new file for writing only
	file, err := os.OpenFile(
		toFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write bytes to file
	bytesWritten, err := file.Write(compressedTextAsByteRay)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}

func getBytesForInt(length int) []byte {

	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, uint64(length))

	return b
}
