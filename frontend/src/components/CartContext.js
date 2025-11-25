import { createContext, useContext, useState, useMemo, useEffect, useCallback } from 'react';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { purchaseApi } from '../services/api';

export const CartContext = createContext();

export const useCart = () => {
  return useContext(CartContext);
};

export const CartProvider = ({ children }) => {
  const [cartItems, setCartItems] = useState([]);
  const { user } = useAuth();

  const fetchCartItems = useCallback(async () => {
    if (!user || user.role !== 'tourist') {
        setCartItems([]);
        return;
    }

    try {
        const response = await purchaseApi.get("");

        if (response.data?.cart?.items) {
            setCartItems(response.data.cart.items);
        } else {
            setCartItems([]);
        }

    } catch (error) {
        console.error("Failed to fetch cart:", error);
        setCartItems([]);
    }
   }, [user]);


    useEffect(() => {
        fetchCartItems();
    }, [fetchCartItems]); 

    const addToCart = (tour) => {
        setCartItems((prevItems) => {
            const exists = prevItems.find((item) => item.tourId === tour.tourId);
            if (exists) {
                return prevItems;
            }
            return [...prevItems, { ...tour }];
        });
    };

    const removeFromCart = async (tourId) => {
        if (!user || user.role !== 'tourist') return;

        try {
            const response = await purchaseApi.delete(`/${tourId}`);
            
            if (response.status === 200) {
                setCartItems((prevItems) => prevItems.filter((item) => item.tourId !== tourId));
            }
        } catch (error) {
            console.error('Greška pri uklanjanju ture iz korpe:', error);
            alert(`Greška pri uklanjanju ture: ${error.response?.data?.error || 'Pokušajte ponovo.'}`);
        }
    };

    const cartItemCount = useMemo(() => cartItems.reduce((count, item) => count + item.quantity, 0), [cartItems]);

    const cartTotal = useMemo(() => {
        return cartItems.reduce((total, item) => total + (item.price * item.quantity), 0).toFixed(2);
    }, [cartItems]);
    
    const value = {
        cartItems,
        cartItemCount,
        cartTotal,
        addToCart,
        removeFromCart,
        fetchCartItems
    };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};