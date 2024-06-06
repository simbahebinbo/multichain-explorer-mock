package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Blob represents the blob data structure
type Blob struct {
	BlobID      string  `json:"BlobID"`
	Status      string  `json:"Status"`
	Commitment  string  `json:"Commitment"`
	BlockNum       int   `json:"Block"`
	Timestamp   string  `json:"Timestamp"`
	Fee         float64 `json:"Fee"`
	Validator   string  `json:"Validator"`
	TxHash      string   `json:"TxHash"`
	State       string   `json:"State"`
}

// Node represents the node data structure
type Node struct {
	NodeString  string `json:"node_string"`
	NodeAddress string `json:"node_address"`
	Chain       string `json:"chain"`
	NodeType    string `json:"node_type"`
}

type BlobDetail struct {
	BlobID       string  `json:"BlobID"`
	Status       string  `json:"Status"`
	Commitment   string  `json:"Commitment"`
	BlockNum        int     `json:"BlockNum"`
	Timestamp    string  `json:"Timestamp"`
	Fee          float64 `json:"Fee"`
	Validator    string  `json:"Validator"`
	Size         int     `json:"Size,omitempty"`
	StorageState string  `json:"StorageState,omitempty"`
	CommitmentXY struct {
		X string `json:"x"`
		Y string `json:"y"`
	}     `json:"Commitment_xy,omitempty"`
	Proof        string  `json:"Proof,omitempty"`
	Data         string  `json:"Data,omitempty"`
}

// Validator represents the validator data structure
type Validator struct {
	ValidatorName           string  `json:"validator_name"`
	ValidatorAddress        string  `json:"validator_address"`
	ValidatorStatus         string  `json:"validator_status"`
	TotalStakedAmount       float64 `json:"total_staked_amount"`
	AvailableStakedAmount   float64 `json:"available_staked_amount"`
	CommissionRate          float64 `json:"commission_rate"`
	VotingPower             float64 `json:"voting_power"`
}


// Sample data for demonstration
var blobs = []Blob{
	{"1", "Confirmed", "Commit1", 100, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000001","Valid"},
	{"2", "Pending", "Commit2", 101, "2024-06-02T12:00:00Z", 0.02, "Validator2","0x0000002","Valid"},
	{"3", "Failed", "Commit3", 102, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000003","Valid"},
	{"4", "Confirmed", "Commit4", 103, "2024-06-01T12:00:00Z", 0.01, "Validator3","0x0000004","Valid"},
	{"5", "Confirmed", "Commit5", 104, "2024-06-01T12:00:00Z", 0.01, "Validator4","0x0000005","Valid"},
	{"6", "Confirmed", "Commit6", 105, "2024-06-01T12:00:00Z", 0.01, "Validator5","0x0000006","Valid"},
}

var btc_blobs = []Blob{
	{"1", "Confirmed", "Commit1", 100, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000001","Valid"},
	{"2", "Pending", "Commit2", 101, "2024-06-02T12:00:00Z", 0.02, "Validator2","0x0000002","Valid"},
	{"3", "Failed", "Commit3", 102, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000003","Valid"},
	{"4", "Confirmed", "Commit4", 103, "2024-06-01T12:00:00Z", 0.01, "Validator3","0x0000004","Valid"},
	{"5", "Confirmed", "Commit5", 104, "2024-06-01T12:00:00Z", 0.01, "Validator4","0x0000005","Valid"},
	{"6", "Confirmed", "Commit6", 105, "2024-06-01T12:00:00Z", 0.01, "Validator5","0x0000006","Valid"},
	{"7", "Confirmed", "Commit7", 106, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000007","Valid"},
	{"8", "Confirmed", "Commit8", 107, "2024-06-01T12:00:00Z", 0.01, "Validator5","0x0000008","Valid"},
	{"9", "Confirmed", "Commit9", 108, "2024-06-01T12:00:00Z", 0.01, "Validator2","0x0000009","Valid"},
	{"10", "Confirmed", "Commit10", 109, "2024-06-01T12:00:00Z", 0.01, "Validator6","0x0000010","Valid"},
	{"11", "Confirmed", "Commit11", 110, "2024-06-01T12:00:00Z", 0.01, "Validator3","0x0000011","Valid"},
	{"12", "Confirmed", "Commit12", 111, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000012","Valid"},
	{"13", "Confirmed", "Commit13", 112, "2024-06-01T12:00:00Z", 0.01, "Validator2","0x0000013","Valid"},
	{"14", "Confirmed", "Commit14", 113, "2024-06-01T12:00:00Z", 0.01, "Validator7","0x0000014","Valid"},
	{"15", "Confirmed", "Commit15", 114, "2024-06-01T12:00:00Z", 0.01, "Validator7","0x0000015","Valid"},
	{"16", "Confirmed", "Commit16", 115, "2024-06-01T12:00:00Z", 0.01, "Validator7","0x0000016","Inalid"},
	{"17", "Confirmed", "Commit17", 116, "2024-06-01T12:00:00Z", 0.01, "Validator5","0x0000017","Inalid"},
	{"18", "Confirmed", "Commit18", 117, "2024-06-01T12:00:00Z", 0.01, "Validator2","0x0000018","Inalid"},
	{"19", "Confirmed", "Commit19", 118, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000019","Inalid"},
	{"20", "Confirmed", "Commit20", 119, "2024-06-01T12:00:00Z", 0.01, "Validator1","0x0000020","Inalid"},
}

// Sample data for demonstration
var blobDetails = []BlobDetail{
	{"1", "Confirmed", "0x1234567890abcdef", 100, "2024-06-01T12:00:00Z", 0.01, "Validator1", 1024, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "13258099556300711131786106409830610145994596628458885637226012245852998915913",
		Y: "11868554521347503492532980178914472193409060128712507356093850651849176305797",
	}, "0x1234567890abcdef", "https://example.com/image1.jpg"},
	{"2", "Pending", "0xabcdef1234567891", 101, "2024-06-02T12:00:00Z", 0.02, "Validator2", 2048, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "0987654321098765432109876543210987654321098765432109876543210987654321",
		Y: "1234567890123456789012345678901234567890123456789012345678901234567890",
	}, "0xabcdef1234567890", "https://example.com/image2.jpg"},
	{"3", "Confirmed", "0x1234567891abcdef", 102, "2024-06-01T12:00:00Z", 0.01, "Validator1", 1024, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "13258099556300711131786106409830610145994596628458885637226012245852998915913",
		Y: "11868554521347503492532980178914472193409060128712507356093850651849176305797",
	}, "0x1234567890abcdef", "https://example.com/image1.jpg"},
	{"4", "Pending", "0xabcdef1234567892", 103, "2024-06-02T12:00:00Z", 0.02, "Validator3", 2048, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "0987654321098765432109876543210987654321098765432109876543210987654321",
		Y: "1234567890123456789012345678901234567890123456789012345678901234567890",
	}, "0xabcdef1234567890", "https://example.com/image2.jpg"},
	{"5", "Confirmed", "0x1234567892abcdef", 104, "2024-06-01T12:00:00Z", 0.01, "Validator2", 1024, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "13258099556300711131786106409830610145994596628458885637226012245852998915913",
		Y: "11868554521347503492532980178914472193409060128712507356093850651849176305797",
	}, "0x1234567890abcdef", "https://example.com/image1.jpg"},
	{"6", "Pending", "0xabcdef1234567893", 105, "2024-06-02T12:00:00Z", 0.02, "Validator1", 2048, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "0987654321098765432109876543210987654321098765432109876543210987654321",
		Y: "1234567890123456789012345678901234567890123456789012345678901234567890",
	}, "0xabcdef1234567890", "https://example.com/image2.jpg"},
	{"7", "Confirmed", "0x1234567893abcdef", 106, "2024-06-01T12:00:00Z", 0.01, "Validator1", 1024, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "13258099556300711131786106409830610145994596628458885637226012245852998915913",
		Y: "11868554521347503492532980178914472193409060128712507356093850651849176305797",
	}, "0x1234567890abcdef", "https://example.com/image1.jpg"},
	{"8", "Pending", "0xabcdef1234567894", 107, "2024-06-02T12:00:00Z", 0.02, "Validator2", 2048, "valid", struct {
		X string `json:"x"`
		Y string `json:"y"`
	}{
		X: "0987654321098765432109876543210987654321098765432109876543210987654321",
		Y: "1234567890123456789012345678901234567890123456789012345678901234567890",
	}, "0xabcdef1234567890", "https://example.com/image2.jpg"},
}

var nodes = []Node{
	{"Node 1", "0xdjshfdcnvnk324fvf7v78vb89bu98vbv8b", "btc", "Broadcast"},
	{"Node 2", "0xdjshfdcnvnk324fvf7v78vb89bu98vbv8b", "btc", "Storage"},
	{"Node 3", "0xdjshfdcnvnk324fvf7v78vb89bu98vbv11", "eth", "Storage"},
	{"Node 4", "0x123sdsfdcnvnk324fvf7v78v89buvbv812", "eth", "Broadcast"},
}

var validators = []Validator{
	{"Validator1", "0x1234567890abcdef1234567890abcdef", "Active", 1000, 800, 0.1, 50},
	{"Validator2", "0xabcdef1234567890abcdef1234567890", "Inactive", 500, 300, 0.2, 20},
}

// HomeDataHandler handles the GET /api/home-data endpoint
func HomeDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		chain := r.URL.Query().Get("chain")
		if chain == "" {
			http.Error(w, "Missing chain parameter", http.StatusBadRequest)
			return
		}

		switch chain {

		case "btc":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(btc_blobs)
		case "eth":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(blobs)

		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// CreateBlobHandler handles the POST /api/create-blob endpoint
func CreateBlobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newBlob Blob
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &newBlob)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		blobs = append(blobs, newBlob)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBlob)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// SearchHandler handles the GET /api/search endpoint with query parameters
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		query := r.URL.Query().Get("q")
		category := r.URL.Query().Get("category")

		if query == "" || category == "" {
			http.Error(w, "Missing query or category parameter", http.StatusBadRequest)
			return
		}

		var results []Blob
		switch category {
		case "Validator", "TxHash", "BlobID", "Commitment", "BlockNum":
			for _, blob := range btc_blobs {
				if category == "BlobID" && strings.Contains(blob.BlobID, query) ||
					category == "Commitment" && strings.Contains(blob.Commitment, query) ||
					category == "BlockNum" && fmt.Sprintf("%d", blob.BlockNum) == query ||
					category == "TxHash" && fmt.Sprintf("%d", blob.TxHash) == query ||
					category == "Validator" && fmt.Sprintf("%d", blob.Validator) == query{
					results = append(results, blob)
				}
			}
		default:
			http.Error(w, "Invalid category parameter", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// BtcBlobsHandler handles the GET /api/btc-blobs endpoint with pagination and filtering
func BtcBlobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		chain := r.URL.Query().Get("chain")
		if chain != "btc" {
			http.Error(w, "Invalid or missing chain parameter", http.StatusBadRequest)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}

		if chain == "" || pageStr == "" {
			http.Error(w, "Missing chain or pageStr parameter", http.StatusBadRequest)
			return
		}

		filter := r.URL.Query().Get("filter")

		// Filter blobs based on the provided filter parameter
		var filteredBlobs []Blob
		for _, blob := range blobs {
			if filter == "" || strings.Contains(blob.BlobID, filter) || strings.Contains(blob.Commitment, filter) ||
				strings.Contains(blob.Status, filter) || strings.Contains(blob.Validator, filter) ||
				strings.Contains(blob.TxHash, filter) {
				filteredBlobs = append(filteredBlobs, blob)
			}
		}

		// Pagination
		const perPage = 10
		total := len(filteredBlobs)
		start := (page - 1) * perPage
		end := start + perPage
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		paginatedBlobs := filteredBlobs[start:end]

		// Response
		response := struct {
			Data       []Blob `json:"data"`
			Pagination struct {
				Total   int `json:"total"`
				Page    int `json:"page"`
				PerPage int `json:"perPage"`
			} `json:"pagination"`
		}{
			Data: paginatedBlobs,
			Pagination: struct {
				Total   int `json:"total"`
				Page    int `json:"page"`
				PerPage int `json:"perPage"`
			}{
				Total:   total,
				Page:    page,
				PerPage: perPage,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// BtcBlobsHandler handles the GET /api/btc-blobs endpoint with pagination and filtering
func EthBlobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		chain := r.URL.Query().Get("chain")
		if chain != "eth" {
			http.Error(w, "Invalid or missing chain parameter", http.StatusBadRequest)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}

		if chain == "" || pageStr == "" {
			http.Error(w, "Missing chain or pageStr parameter", http.StatusBadRequest)
			return
		}

		filter := r.URL.Query().Get("filter")

		// Filter blobs based on the provided filter parameter
		var filteredBlobs []Blob
		for _, blob := range blobs {
			if filter == "" || strings.Contains(blob.BlobID, filter) || strings.Contains(blob.Commitment, filter) ||
				strings.Contains(blob.Status, filter) || strings.Contains(blob.Validator, filter) ||
				strings.Contains(blob.TxHash, filter) {
				filteredBlobs = append(filteredBlobs, blob)
			}
		}

		// Pagination
		const perPage = 10
		total := len(filteredBlobs)
		start := (page - 1) * perPage
		end := start + perPage
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		paginatedBlobs := filteredBlobs[start:end]

		// Response
		response := struct {
			Data       []Blob `json:"data"`
			Pagination struct {
				Total   int `json:"total"`
				Page    int `json:"page"`
				PerPage int `json:"perPage"`
			} `json:"pagination"`
		}{
			Data: paginatedBlobs,
			Pagination: struct {
				Total   int `json:"total"`
				Page    int `json:"page"`
				PerPage int `json:"perPage"`
			}{
				Total:   total,
				Page:    page,
				PerPage: perPage,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func BlobDetailHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		blobID := vars["blobID"]
		chain := vars["chain"]

		if chain != "btc" && chain != "eth" {
			http.Error(w, "Invalid chain parameter", http.StatusBadRequest)
			return
		}

		var foundBlob *BlobDetail
		for _, blob := range blobDetails {
			if blob.BlobID == blobID {
				foundBlob = &blob
				break
			}
		}

		if foundBlob == nil {
			http.Error(w, "Blob not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(foundBlob)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
// NodesHandler handles the GET /api/nodes endpoint
func NodesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		chain := r.URL.Query().Get("chain")

		var filteredNodes []Node
		if chain == "" {
			filteredNodes = nodes
		} else {
			for _, node := range nodes {
				if node.Chain == chain {
					filteredNodes = append(filteredNodes, node)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(filteredNodes)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// GetValidatorHandler handles the GET /api/getValidator endpoint
func GetValidatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if len(validators) == 0 {
			http.Error(w, "No validators found", http.StatusNotFound)
			return
		}
		validator := validators // Assuming a single validator for simplicity
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(validator)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/home-data", HomeDataHandler)
	r.HandleFunc("/api/create-blob", CreateBlobHandler)
	r.HandleFunc("/api/search", SearchHandler)
	r.HandleFunc("/api/btc-blobs", BtcBlobsHandler)
	r.HandleFunc("/api/eth-blobs", EthBlobsHandler)
	r.HandleFunc("/api/blob-detail", BlobDetailHandler)
	r.HandleFunc("/api/nodes", NodesHandler)
	r.HandleFunc("/api/getValidator", GetValidatorHandler)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
