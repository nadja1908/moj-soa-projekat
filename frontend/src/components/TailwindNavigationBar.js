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
    <nav className="bg-white shadow-md border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-6">
        <div className="flex justify-between items-center h-20">
          {/* LOGO */}
          <Link to="/" className="flex items-center space-x-3 group">
            <img
              src="/logo.png"
              alt="Wander Vistas"
              className="w-40 md:w-52 h-auto transition-transform duration-300 group-hover:scale-105"
            />
          </Link>

          {/* USER SECTION */}
          <div className="flex items-center space-x-4">
            {user ? (
              <div className="flex items-center space-x-3">
                {/* Avatar */}
                <div className="w-10 h-10 bg-gradient-to-br from-[#0d1b2a] to-[#1b263b] rounded-full flex items-center justify-center text-white font-bold text-lg">
                  {user.username.charAt(0).toUpperCase()}
                </div>

                {/* Info */}
                <div className="hidden md:block">
                  <div className="text-sm font-semibold text-gray-900">
                    <Link 
                    className="inline-block hover:text-indigo-600 transition-colors duration-200 cursor-pointer"
                    to="/profile">{user.username}</Link>
                  </div>
                  <div
                    className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-semibold ${getRoleBadgeStyle(
                      user.role
                    )}`}
                  >
                    <span className="mr-1">{getRoleIcon(user.role)}</span>
                    {getRoleDisplayName(user.role)}
                  </div>
                </div>

                {/* Logout */}
                <button
                  onClick={handleLogout}
                  className="bg-gradient-to-r from-[#e76f51] to-[#f4a261] hover:from-[#e85d44] hover:to-[#f1a550] text-white px-4 py-2 rounded-lg font-medium transition-all duration-200 hover:shadow-md flex items-center space-x-1"
                >
                  <span>🚪</span>
                  <span className="hidden md:inline">Odjavi se</span>
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-3">
                <Link
                  to="/login"
                  className="text-[#0d1b2a] hover:text-[#e76f51] font-semibold transition-colors duration-200 flex items-center space-x-1"
                >
                  <span>🔐</span>
                  <span>Prijavi se</span>
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
