import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import NavigationBar from './components/NavigationBar';
import Login from './components/Login';
import Register from './components/Register';
import AdminPanel from './components/AdminPanel';
import TourManagement from './components/TourManagement';
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
    switch (user.role) {
      case 'administrator': return "/admin";
      case 'guide': return "/guide/tours";
      case 'tourist': return "/tours";
      default: return "/login";
    }
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
            element={user && user.role === 'tourist' ? <div>Tourist Tours View (TODO)</div> : <Navigate to="/login" />} 
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