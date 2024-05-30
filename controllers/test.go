package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
    

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetFacID(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		log.Print("ID :", id)

		// Search By ID
		// objID, err := primitive.ObjectIDFromHex(id)
		// if err != nil {
		// 	http.Error(w, "Invalid ID format", http.StatusBadRequest)
		// 	return
		// }
		// filter := bson.M{"facultyid": objID}
		// err = collection.FindOne(context.Background(), filter).Decode(&fac)
		// var fac bson.M

		// Search By Anothor Fildes
		var results []bson.M
		filter := bson.M{
			"$or": []bson.M{
				{"_id": id},
				{"product_name": id},
				// เพิ่มฟิลด์อื่นๆ ที่ต้องการค้นหาได้ที่นี่
			},
		}
		// ค้นหาเอกสารที่ตรงกับฟิลเตอร์
		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())
		log.Print(filter)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Document not found", http.StatusNotFound)
				return
			}
		}

		for cursor.Next(context.Background()) {
			var fac bson.M
			if err := cursor.Decode(&fac); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, fac)
		}
		json.NewEncoder(w).Encode(results)
	}
}
