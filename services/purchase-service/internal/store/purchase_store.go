package store

import (
	"context"
	"fmt"
	"time"

	"purchase-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func calculateTotalPrice(items []model.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func (s *Store) GetCartByTouristID(ctx context.Context, touristID int64) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart

	err := s.CartsCollection.FindOne(ctx, bson.M{"touristId": touristID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart, nil
}

func (s *Store) SaveCartWithItem(ctx context.Context, touristID int64, newItem model.OrderItem) (*model.ShoppingCart, error) {
	cart, err := s.GetCartByTouristID(ctx, touristID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = &model.ShoppingCart{
			TouristID: touristID,
			Items:     []model.OrderItem{newItem},
			CreatedAt: time.Now(),
		}
	} else {
		updated := false
		for i, item := range cart.Items {
			if item.TourID == newItem.TourID {
				cart.Items[i].Quantity += newItem.Quantity
				updated = true
				break
			}
		}
		if !updated {
			cart.Items = append(cart.Items, newItem)
		}
	}

	cart.TotalPrice = calculateTotalPrice(cart.Items)
	cart.LastUpdated = time.Now()

	filter := bson.M{"touristId": touristID}
	update := bson.M{"$set": cart}
	opts := options.Update().SetUpsert(true)

	_, err = s.CartsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to save cart: %w", err)
	}

	return s.GetCartByTouristID(ctx, touristID)
}

func (s *Store) RemoveItemFromCart(ctx context.Context, touristID int64, tourID int64) (*model.ShoppingCart, error) {
	cart, err := s.GetCartByTouristID(ctx, touristID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, fmt.Errorf("shopping cart not found")
	}

	newItems := []model.OrderItem{}
	removed := false

	for _, item := range cart.Items {
		if item.TourID != tourID {
			newItems = append(newItems, item)
		} else {
			removed = true
		}
	}

	if !removed {
		return cart, fmt.Errorf("item not found in cart")
	}

	cart.Items = newItems
	cart.TotalPrice = calculateTotalPrice(cart.Items)
	cart.LastUpdated = time.Now()

	_, err = s.CartsCollection.UpdateOne(
		ctx,
		bson.M{"touristId": touristID},
		bson.M{"$set": cart},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return cart, nil
}

func (s *Store) SaveTokens(ctx context.Context, tokens []model.TourPurchaseToken) error {
	if len(tokens) == 0 {
		return nil
	}

	docs := make([]interface{}, len(tokens))
	for i, t := range tokens {
		docs[i] = t
	}

	_, err := s.TokensCollection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	return nil
}

func (s *Store) ClearCart(ctx context.Context, touristID int64) error {
	_, err := s.CartsCollection.DeleteOne(ctx, bson.M{"touristId": touristID})
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}
