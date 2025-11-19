import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const NavigationBar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const getRoleBadgeStyle = (role) => {
    switch (role) {
      case 'administrator':
        return 'bg-red-100 text-red-800 border border-red-200';
      case 'guide':
        return 'bg-green-100 text-green-800 border border-green-200';
      case 'tourist':
        return 'bg-blue-100 text-blue-800 border border-blue-200';
      default:
        return 'bg-gray-100 text-gray-800 border border-gray-200';
    }
  };

  const getRoleDisplayName = (role) => {
    switch (role) {
      case 'administrator':
        return 'Administrator';
      case 'guide':
        return 'Vodič';
      case 'tourist':
        return 'Turista';
      default:
        return role;
    }
  };

  const getRoleIcon = (role) => {
    switch (role) {
      case 'administrator':
        return '👑';
      case 'guide':
        return '🗺️';
      case 'tourist':
        return '🧳';
      default:
        return '👤';
    }
  };

  return (
    <nav className="bg-white shadow-lg border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center space-x-2 text-xl font-bold text-indigo-600 hover:text-indigo-800 transition-colors duration-200">
            <span className="text-2xl">📝</span>
            <span>SOA Blog Platform</span>
          </Link>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center space-x-6">
            {user && (
              <>
                {(user.role === 'guide' || user.role === 'tourist') && (
                  <Link 
                    to="/posts" 
                    className="flex items-center space-x-1 text-gray-700 hover:text-indigo-600 transition-colors duration-200 font-medium"
                  >
                    <span>📰</span>
                    <span>Blog</span>
                  </Link>
                )}
                
                <Link 
                  to="/create-post" 
                  className="flex items-center space-x-1 text-gray-700 hover:text-indigo-600 transition-colors duration-200 font-medium"
                >
                  <span>✍️</span>
                  <span>Novi Post</span>
                </Link>
                
                {user.role === 'administrator' && (
                  <Link 
                    to="/admin" 
                    className="flex items-center space-x-1 text-gray-700 hover:text-indigo-600 transition-colors duration-200 font-medium"
                  >
                    <span>👑</span>
                    <span>Admin Panel</span>
                  </Link>
                )}
              </>
            )}
          </div>

          {/* User Section */}
          <div className="flex items-center space-x-4">
            {user ? (
              <div className="flex items-center space-x-3">
                {/* User Avatar */}
                <div className="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-full flex items-center justify-center text-white font-bold text-lg">
                  {user.username.charAt(0).toUpperCase()}
                </div>
                
                {/* User Info */}
                <div className="hidden md:block">
                  <div className="text-sm font-semibold text-gray-900">{user.username}</div>
                  <div className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-semibold ${getRoleBadgeStyle(user.role)}`}>
                    <span className="mr-1">{getRoleIcon(user.role)}</span>
                    {getRoleDisplayName(user.role)}
                  </div>
                </div>
                
                {/* Logout Button */}
                <button 
                  onClick={handleLogout}
                  className="bg-red-50 hover:bg-red-100 text-red-700 border border-red-200 px-4 py-2 rounded-lg font-medium transition-all duration-200 hover:shadow-md flex items-center space-x-1"
                >
                  <span>🚪</span>
                  <span className="hidden md:inline">Odjavi se</span>
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-3">
                <Link 
                  to="/login" 
                  className="text-gray-700 hover:text-indigo-600 font-medium transition-colors duration-200 flex items-center space-x-1"
                >
                  <span>🔐</span>
                  <span>Prijavi se</span>
                </Link>
                <Link 
                  to="/register" 
                  className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-medium transition-all duration-200 hover:shadow-md flex items-center space-x-1"
                >
                  <span>📝</span>
                  <span>Registruj se</span>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default NavigationBar;