import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import NavigationBar from './components/NavigationBar';
import Login from './components/Login';
import Register from './components/Register';
import AdminPanel from './components/AdminPanel';
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

  return (
    <div className="min-h-screen bg-white">
      <NavigationBar />
      <div className="main-container bg-white">
        <Routes>
          <Route 
            path="/login" 
            element={!user ? <Login /> : <Navigate to={user.role === 'administrator' ? "/admin" : "/login"} />} 
          />
          <Route 
            path="/register" 
            element={!user ? <Register /> : <Navigate to={user.role === 'administrator' ? "/admin" : "/login"} />} 
          />
          <Route 
            path="/admin" 
            element={user && user.role === 'administrator' ? <AdminPanel /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/" 
            element={<Navigate to={user ? (user.role === 'administrator' ? "/admin" : "/login") : "/login"} />} 
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