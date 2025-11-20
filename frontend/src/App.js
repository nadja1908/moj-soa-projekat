import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import NavigationBar from './components/NavigationBar';
import Login from './components/Login';
import Register from './components/Register';
import AdminPanel from './components/AdminPanel';
import TourManagement from './components/TourManagement';
import TouristTours from './components/TouristTours';
import BlogPosts from './components/BlogPosts';
import CreatePost from './components/CreatePost';
import { AuthProvider, useAuth } from './context/AuthContext';

function AppContent() {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center" style={{ height: '100vh' }}>
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  const getRedirectPath = (user) => {
    if (!user) return "/login";
    // Svi korisnici idu na blog kao početnu stranicu
    return "/posts";
  };

  return (
    <div className="min-h-screen bg-white">
      <NavigationBar />
      <div className="main-container bg-white">
        <Routes>
          <Route 
            path="/login" 
            element={!user ? <Login /> : <Navigate to={getRedirectPath(user)} />} 
          />
          <Route 
            path="/register" 
            element={!user ? <Register /> : <Navigate to={getRedirectPath(user)} />} 
          />
          
          {/* Administrator rute */}
          <Route 
            path="/admin" 
            element={user && user.role === 'administrator' ? <AdminPanel /> : <Navigate to="/login" />} 
          />
          
          {/* Guide rute */}
          <Route 
            path="/guide/tours" 
            element={user && user.role === 'guide' ? <TourManagement /> : <Navigate to="/login" />} 
          />
          
          {/* Tourist rute */}
          <Route 
            path="/tours" 
            element={user && user.role === 'tourist' ? <TouristTours /> : <Navigate to="/login" />} 
          />
          
          {/* Blog rute - dostupne svim korisnicima */}
          <Route 
            path="/posts" 
            element={<BlogPosts />} 
          />
          
          {/* Create Post - za sve autentifikovane korisnike */}
          <Route 
            path="/create-post" 
            element={user ? <CreatePost /> : <Navigate to="/login" />} 
          />
          
          <Route 
            path="/" 
            element={<Navigate to={getRedirectPath(user)} />} 
          />
          
          {/* Catch-all route za sve nepoznate putanje */}
          <Route 
            path="*" 
            element={<Navigate to={getRedirectPath(user)} />} 
          />
        </Routes>
      </div>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;