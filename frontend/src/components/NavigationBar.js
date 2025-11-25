import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import TravelLogo from './TravelLogo';
import { useCart } from './CartContext';

const NavigationBar = () => {
  const { user, logout } = useAuth();
  const { cartItemCount } = useCart();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const getRoleBadgeStyle = (role) => {
    switch (role) {
      case 'administrator':
        return 'bg-gradient-to-r from-orange-100 to-amber-100 text-orange-800 border border-orange-200';
      case 'guide':
        return 'bg-gradient-to-r from-amber-100 to-yellow-100 text-amber-800 border border-amber-200';
      case 'tourist':
        return 'bg-gradient-to-r from-rose-100 to-orange-100 text-rose-800 border border-orange-200';
      default:
        return 'bg-gradient-to-r from-gray-100 to-slate-100 text-gray-800 border border-gray-200';
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

  const getNavigationLinks = () => {
    if (!user) {
      return null;
    }
    
    // Bazni linkovi za sve autentifikovane korisnike
    const baseLinks = null;
    
    // Blog link samo za guide i tourist (ne za admin)
    const blogLink = (user.role === 'guide' || user.role === 'tourist') ? (
      <>
        <Link
          to="/posts"
          className="text-gray-700 hover:text-orange-600 font-medium transition-colors duration-200 text-sm xl:text-base"
        >
          📰 Blog
        </Link>
        <Link
          to="/my-posts"
          className="text-gray-700 hover:text-orange-600 font-medium transition-colors duration-200 text-sm xl:text-base"
        >
          📝 Moji Blogovi
        </Link>
        <Link
          to="/community"
          className="text-gray-700 hover:text-purple-600 font-medium transition-colors duration-200 text-sm xl:text-base"
        >
          👥 Zajednica
        </Link>
      </>
    ) : null;
    
    const roleSpecificLinks = (() => {
      switch (user.role) {
        case 'administrator':
          return (
            <>
              <Link
                to="/create-post"
                className="text-gray-700 hover:text-amber-600 font-medium text-sm xl:text-base transition-colors duration-200"
              >
                ✍️ Novi Post
              </Link>
              <Link
                to="/admin"
                className="text-gray-700 hover:text-rose-600 font-medium text-sm xl:text-base transition-colors duration-200"
              >
                👑 Admin Panel
              </Link>
            </>
          );
        case 'guide':
          return (
            <>
              <Link
                to="/create-post"
                className="text-gray-700 hover:text-amber-600 font-medium text-sm xl:text-base transition-colors duration-200"
              >
                ✍️ Novi Post
              </Link>
              <Link
                to="/guide/tours"
                className="text-gray-700 hover:text-orange-600 font-medium transition-colors duration-200 text-sm xl:text-base flex items-center"
              >
                🗺️ Moje ture
              </Link>
            </>
          );
        case 'tourist':
          return (
            <>
              <Link
                to="/create-post"
                className="text-gray-700 hover:text-amber-600 font-medium text-sm xl:text-base transition-colors duration-200"
              >
                ✍️ Novi Post
              </Link>
              <Link
                to="/tours"
                className="text-gray-700 hover:text-orange-600 font-medium transition-colors duration-200 text-sm xl:text-base flex items-center"
              >
                🎒 Dostupne ture
              </Link>
              <Link
                to="/purchase"
                className="relative text-gray-700 hover:text-orange-600 font-medium transition-colors duration-200 text-sm xl:text-base flex items-center"
                title="Moja korpa"
              >
                🛒 Korpa
                {cartItemCount > 0 && (
                  <span className="absolute top-0 right-0 transform translate-x-1/2 -translate-y-1/2 bg-rose-500 text-white text-xs font-bold px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
                    {cartItemCount}
                  </span>
                )}
              </Link>
            </>
          );
        default:
          return null;
      }
    })();
    
    return (
      <>
        {blogLink}
        {roleSpecificLinks}
      </>
    );
  };

  return (
    <nav className="bg-white shadow-lg border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-6xl xl:max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-24 xl:h-28">
          
          {/* LOGO */}
          <Link
            to="/"
            className="flex items-center shrink-0"
          >
            <TravelLogo className="w-40 h-20 xl:w-48 xl:h-24" />
          </Link>

          {/* Desktop meni */}
          <div className="hidden md:flex items-center space-x-4 xl:space-x-6">
            {getNavigationLinks()}
          </div>

          {/* User sekcija */}
          <div className="flex items-center space-x-3 xl:space-x-4">
            {user ? (
              <div className="flex items-center space-x-3 xl:space-x-4">
                {/* Avatar */}
                <div className="w-9 xl:w-10 h-9 xl:h-10 bg-gradient-to-br from-orange-400 via-amber-500 to-rose-500 rounded-full flex items-center justify-center text-white font-bold text-sm xl:text-lg">
                  {user.username.charAt(0).toUpperCase()}
                </div>

                {/* Info */}
                <div className="hidden md:block">
                  <div className="text-sm font-semibold text-gray-900">
                    <Link to="/profile"
                          className="inline-block hover:text-indigo-600 transition-colors duration-200 cursor-pointer"
                    >{user.username}</Link>
                  </div>
                  <div
                    className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-semibold ${getRoleBadgeStyle(
                      user.role
                    )}`}
                  >
                    {getRoleDisplayName(user.role)}
                  </div>
                </div>

                {/* Logout */}
                <button
                  onClick={handleLogout}
                  className="bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white px-3 xl:px-4 py-2 rounded-lg font-medium transition-all duration-200 hover:shadow-md text-sm xl:text-base"
                >
                  Odjavi se
                </button>
              </div>
            ) : (
              <div className="flex items-center space-x-2 xl:space-x-3">
                <Link
                  to="/login"
                  className="text-gray-700 hover:text-orange-600 font-medium transition-colors duration-200 text-sm xl:text-base"
                >
                  Prijavi se
                </Link>
                <Link
                  to="/register"
                  className="bg-gradient-to-r from-[#e76f51] to-[#f4a261] hover:from-[#e85d44] hover:to-[#f1a550] text-white px-3 xl:px-4 py-2 rounded-lg font-medium transition-all duration-200 hover:shadow-md text-sm xl:text-base"
                >
                  Registruj se
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
