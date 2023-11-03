// //go:build integration
// // +build integration
package weaviateclient

import (
	"fmt"
	"lucidify-api/data/store/storemodels"
	"testing"

	"github.com/google/uuid"
)

func TestUploadDeleteChunk(t *testing.T) {
	weaviateClient, err := NewWeaviateClientTest()
	if err != nil {
		t.Errorf("failed to create weaviate client: %v", err)
	}

	documentID := uuid.New()
	// Create a sample DocumentChunk
	chunk := storemodels.Chunk{
		ChunkID:      uuid.New(),
		UserID:       uuid.New().String(),
		DocumentID:   documentID,
		ChunkContent: "Test chunk content",
		ChunkIndex:   0,
	}

	err = weaviateClient.UploadChunk(chunk)
	if err != nil {
		t.Errorf("UploadChunk failed: %v", err)
	}
	err = weaviateClient.UploadChunk(chunk)
	if err == nil {
		t.Errorf("UploadChunk should have failed due to duplication: %v", err)
	}
	err = weaviateClient.DeleteChunk(chunk.ChunkID)
	if err != nil {
		t.Errorf("DeleteAllChunksByDocumentID failed: %v", err)
	}
	err = weaviateClient.UploadChunk(chunk)
	if err != nil {
		t.Errorf("re-UploadChunk failed proceeding DeleteChunk: %v", err)
	}
	err = weaviateClient.DeleteChunk(chunk.ChunkID)
	if err != nil {
		t.Errorf("DeleteAllChunksByDocumentID failed: %v", err)
	}
}

func TestUploadDeleteChunks(t *testing.T) {
	weaviateClient, err := NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	documentID := uuid.New()
	userID := uuid.New().String()
	var chunks []storemodels.Chunk

	// Create a sample DocumentChunk
	chunk0 := storemodels.Chunk{
		ChunkID:      uuid.New(),
		UserID:       userID,
		DocumentID:   documentID,
		ChunkContent: "Test chunk content",
		ChunkIndex:   0,
	}
	chunks = append(chunks, chunk0)

	// Create a sample DocumentChunk
	chunk1 := storemodels.Chunk{
		ChunkID:      uuid.New(),
		UserID:       userID,
		DocumentID:   documentID,
		ChunkContent: "Test chunk content",
		ChunkIndex:   1,
	}
	chunks = append(chunks, chunk1)

	err = weaviateClient.UploadChunks(chunks)
	if err != nil {
		t.Errorf("UploadChunks failed: %v", err)
	}
	err = weaviateClient.UploadChunks(chunks)
	if err == nil {
		t.Errorf("UploadChunks should have failed due to duplication: %v", err)
	}
	err = weaviateClient.DeleteChunks(chunks)
	if err != nil {
		t.Errorf("DeleteAllChunksByDocumentID failed: %v", err)
	}
	err = weaviateClient.UploadChunks(chunks)
	if err != nil {
		t.Errorf("re-UploadChunks failed proceeding DeleteChunk: %v", err)
	}
}

func getTestChunks() []storemodels.Chunk {
	documentID := uuid.New()
	userID := uuid.New()
	var chunks []storemodels.Chunk
	chunk0 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     userID.String(),
		DocumentID: documentID,
		ChunkContent: "Cats, with their graceful movements and independent nature, have been revered" +
			" and adored by many civilizations throughout history. In ancient Egypt, they" +
			" were considered sacred and were even associated with the goddess Bastet, who" +
			" was depicted as a lioness or a woman with the head of a lioness. Cats were" +
			" believed to have protective qualities, and harming one was considered a grave" +
			" offense. Their sleek and mysterious demeanor earned them a special place in the" +
			" hearts of the Egyptians, a sentiment that has persisted to modern times.",
		ChunkIndex: 0,
	}
	chunks = append(chunks, chunk0)

	chunk1 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     userID.String(),
		DocumentID: documentID,
		ChunkContent: "One of the most captivating features of a cat is its eyes. With large, round" +
			" pupils that can expand and contract based on the amount of light, a cat's eyes" +
			" are a marvel of evolution. This feature allows them to have excellent night" +
			" vision, making them adept hunters even in low-light conditions. The reflective" +
			" layer behind their retinas, known as the tapetum lucidum, gives their eyes a" +
			" distinctive glow in the dark and further enhances their ability to see at" +
			" night.",
		ChunkIndex: 1,
	}
	chunks = append(chunks, chunk1)

	chunk2 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     userID.String(),
		DocumentID: documentID,
		ChunkContent: "Cats are known for their grooming habits. They spend a significant amount of" +
			" time each day cleaning their fur with their rough tongues. This not only keeps" +
			" them clean but also helps regulate their body temperature. The act of grooming" +
			" also has a calming effect on cats, and it's not uncommon to see them grooming" +
			" themselves or other cats as a sign of affection and bonding. This meticulous" +
			" cleaning ritual also aids in reducing scent, making them stealthier hunters.",
		ChunkIndex: 2,
	}
	chunks = append(chunks, chunk2)

	chunk3 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     userID.String(),
		DocumentID: documentID,
		ChunkContent: "The purring of a cat is a sound that many find soothing and comforting. While" +
			" it's commonly associated with contentment, cats also purr when they are in" +
			" pain, anxious, or even when they're near death. The exact mechanism and purpose" +
			" of purring remain a subject of research and speculation. Some theories suggest" +
			" that purring has healing properties, as the vibrations can stimulate the" +
			" production of certain growth factors that aid in wound healing.",
		ChunkIndex: 3,
	}
	chunks = append(chunks, chunk3)

	chunk4 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     userID.String(),
		DocumentID: documentID,
		ChunkContent: "Domestic cats, despite being pampered pets in many households, still retain" +
			" many of their wild instincts. Their tendency to \"hunt\" toys, pounce on moving" +
			" objects, or even their habit of bringing back prey to their owners are all" +
			" remnants of their wild ancestry. These behaviors are deeply ingrained and serve" +
			" as a reminder that beneath their cuddly exterior lies a skilled predator, honed" +
			" by millions of years of evolution.",
		ChunkIndex: 4,
	}
	chunks = append(chunks, chunk4)

	secondUserID := uuid.New()
	secondDocumentID := uuid.New()

	chunk5 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     secondUserID.String(),
		DocumentID: secondDocumentID,
		ChunkContent: "Introduction to Dogs: Dogs, often referred to as \"man's best friend,\" have been" +
			" companions to humans for thousands of years. Originating from wild wolves," +
			" these loyal creatures have been domesticated and bred for various roles" +
			" throughout history, from hunting and herding to companionship. Their keen" +
			" senses, especially their sense of smell, combined with their innate" +
			" intelligence, make them invaluable partners in numerous tasks. Today, dogs are" +
			" found in countless households worldwide, providing joy, comfort, and sometimes" +
			" even protection to their human families.",
		ChunkIndex: 0,
	}
	chunks = append(chunks, chunk5)

	chunk6 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     secondUserID.String(),
		DocumentID: secondDocumentID,
		ChunkContent: "Diverse Breeds: The world of dogs is incredibly diverse, with over 340" +
			" recognized breeds, each with its unique characteristics, temperament, and" +
			" appearance. From the tiny Chihuahua to the majestic Great Dane, dogs come in" +
			" all shapes and sizes. Some breeds, like the Border Collie, are known for their" +
			" intelligence and agility, while others, such as the Saint Bernard, are" +
			" celebrated for their strength and gentle nature. This vast diversity ensures" +
			" that there's a perfect dog breed for almost every individual and lifestyle.",
		ChunkIndex: 1,
	}
	chunks = append(chunks, chunk6)

	chunk7 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     secondUserID.String(),
		DocumentID: secondDocumentID,
		ChunkContent: "Roles and Responsibilities: Beyond being mere pets, dogs play various roles in" +
			" human societies. Service dogs assist individuals with disabilities, guiding the" +
			" visually impaired or alerting those with hearing loss. Therapy dogs provide" +
			" emotional support in hospitals, schools, and nursing homes, offering comfort to" +
			" those in need. Working dogs, like police K9 units or search and rescue teams," +
			" perform critical tasks that save lives. However, with these roles comes the" +
			" responsibility for owners to provide proper training, care, and attention to" +
			" their canine companions.",
		ChunkIndex: 2,
	}
	chunks = append(chunks, chunk7)

	chunk8 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     secondUserID.String(),
		DocumentID: secondDocumentID,
		ChunkContent: "Health and Care: Just like humans, dogs have specific health and care needs" +
			" that owners must address. Regular veterinary check-ups, vaccinations, and a" +
			" balanced diet are essential for a dog's well-being. Grooming, depending on the" +
			" breed, can range from daily brushing to occasional baths. Exercise is crucial" +
			" for a dog's physical and mental health, with daily walks and playtime being" +
			" beneficial. Additionally, training and socialization from a young age ensure" +
			" that dogs are well-behaved and can interact positively with other animals and" +
			" people.",
		ChunkIndex: 3,
	}
	chunks = append(chunks, chunk8)

	chunk9 := storemodels.Chunk{
		ChunkID:    uuid.New(),
		UserID:     secondUserID.String(),
		DocumentID: secondDocumentID,
		ChunkContent: "The Bond Between Humans and Dogs: The relationship between humans and dogs is" +
			" profound and multifaceted. Dogs offer unconditional love, loyalty, and" +
			" companionship, often becoming integral members of the family. In return, humans" +
			" provide care, shelter, and affection. Numerous studies have shown that owning a" +
			" dog can reduce stress, increase physical activity, and bring joy to their" +
			" owners. This symbiotic relationship, built on mutual trust and respect," +
			" showcases the incredible bond that has existed between our two species for" +
			" millennia.",
		ChunkIndex: 4,
	}
	chunks = append(chunks, chunk9)

	return chunks
}

func TestSearchDocumentsByText(t *testing.T) {
	weaviateClient, err := NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	fmt.Println("Woop woop that's the sound of the beez")

	// Keep track of uploaded document IDs for cleanup
	testChunks := getTestChunks()
	if err != nil {
		t.Errorf("setup failed: %v", err)
	}
	err = weaviateClient.UploadChunks(testChunks)
	if err != nil {
		t.Errorf("UploadChunks failed: %v", err)
	}

	defer func() {
		if err := weaviateClient.DeleteChunks(testChunks); err != nil {
			t.Errorf("teardown failed: %v", err)
		}
	}()
	// Define a query and limit for the test
	top_k := 3
	userID := testChunks[0].UserID

	concepts := []string{"small animal that goes meow sometimes"}

	result, err := weaviateClient.SearchDocumentsByText(top_k, userID, concepts)
	if err != nil {
		t.Errorf("SearchDocumentsByText failed: %v", err)
	}

	if len(result) != top_k {
		t.Errorf("incorrect number of results: got %v, want %v", len(result), top_k)
	}

	for _, chunk := range result {
		fmt.Printf("Chunk: %+v\n", chunk)
	}

	secondUserID := testChunks[5].UserID
	concepts = []string{"small animal that goes meow sometimes"}

	result, err = weaviateClient.SearchDocumentsByText(top_k, secondUserID, concepts)
	if err != nil {
		t.Errorf("SearchDocumentsByText failed: %v", err)
	}

	for _, chunk := range result {
		if chunk.UserID != secondUserID {
			t.Errorf("Should not have returned user1 chunks")
		}
	}
	for _, chunk := range result {
		fmt.Printf("Chunk: %+v\n", chunk)
	}
}

func TestDeleteAllChunksByUserID(t *testing.T) {
	weaviateClient, err := NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	documentID := uuid.New()
	userID := uuid.New().String()
	var chunks []storemodels.Chunk

	// Create a sample DocumentChunk
	chunk0 := storemodels.Chunk{
		ChunkID:      uuid.New(),
		UserID:       userID,
		DocumentID:   documentID,
		ChunkContent: "Test chunk content",
		ChunkIndex:   0,
	}
	chunks = append(chunks, chunk0)

	// Create a sample DocumentChunk
	chunk1 := storemodels.Chunk{
		ChunkID:      uuid.New(),
		UserID:       userID,
		DocumentID:   documentID,
		ChunkContent: "Test chunk content",
		ChunkIndex:   1,
	}
	chunks = append(chunks, chunk1)

	chunk2 := storemodels.Chunk{
		ChunkID:      uuid.New(),
		UserID:       uuid.New().String(),
		DocumentID:   documentID,
		ChunkContent: "Test chunk content",
		ChunkIndex:   2,
	}
	chunks = append(chunks, chunk2)

	err = weaviateClient.UploadChunks(chunks)
	if err != nil {
		t.Errorf("UploadChunks failed: %v", err)
	}
	err = weaviateClient.UploadChunks(chunks)
	if err == nil {
		t.Errorf("UploadChunks should have failed due to duplication: %v", err)
	}
	chunks, err = weaviateClient.GetChunks(chunks)
	if err != nil || len(chunks) != 3 {
		t.Errorf("GetChunks should have succeeded: %v", err)
	}
	err = weaviateClient.DeleteAllChunksByUserID(userID)
	if err != nil {
		t.Errorf("DeleteAllChunksByDocumentID failed: %v", err)
	}
	chunks, err = weaviateClient.GetChunks(chunks)
	if err == nil || len(chunks) != 1 {
		t.Errorf("Get Chunks should return 1 chunk. returned chunks: %v", len(chunks))
	}
}
