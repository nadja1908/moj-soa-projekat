import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Login = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const result = await login(formData.username, formData.password);
    if (result.success) {
      navigate('/dashboard');
    } else {
      setError(result.error);
    }
    setLoading(false);
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  return (
    <div className="min-h-[calc(100vh-6rem)] flex bg-white overflow-hidden">
      {/* Leva slika */}
      <div className="hidden md:block md:w-1/2 relative">
        <img
          src="/newyork.jpg"
          alt="Travel destinations"
          className="w-full h-full object-cover brightness-105 contrast-105 saturate-110"
        />
        <div className="absolute inset-0 bg-gradient-to-r from-orange-100/25 via-transparent to-transparent pointer-events-none" />
      </div>

      {/* Desna strana */}
      <div className="w-full md:w-1/2 flex items-center justify-center px-4">
        <div className="bg-white rounded-2xl shadow-lg w-full max-w-sm border border-orange-50">
          
          {/* Header */}
          <div className="px-6 pt-4 pb-2 text-center">
            <h2 className="text-base font-semibold text-[#222]">
              Dobrodošli nazad
            </h2>
            <p className="text-[11px] text-[#7a7a7a]">
              Prijavite se i nastavite svoja putovanja sa nama
            </p>
          </div>

          {/* Forma */}
          <div className="px-6 pb-4">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-3 py-1.5 rounded-lg mb-2 flex items-center text-[11px]">
                <svg className="w-3.5 h-3.5 mr-1.5" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-3">
              {/* Username */}
              <div>
                <label className="block text-[11px] font-medium text-gray-700 mb-1">
                  Korisničko ime
                </label>
                <div className="relative">
                  <svg
                    className="absolute left-3 top-2.5 w-4 h-4 text-gray-400 pointer-events-none"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                      clipRule="evenodd"
                    />
                  </svg>
                  <input
                    type="text"
                    name="username"
                    value={formData.username}
                    onChange={handleChange}
                    placeholder="Unesite korisničko ime"
                    required
                    className="w-full px-10 py-2 bg-[#fdf4ec] border border-orange-100 rounded-lg focus:ring-2 focus:ring-[#f4a261] focus:border-[#f4a261] text-[13px] text-gray-900"
                  />
                </div>
              </div>

              {/* Password */}
              <div>
                <label className="block text-[11px] font-medium text-gray-700 mb-1">
                  Lozinka
                </label>
                <div className="relative">
                  <svg
                    className="absolute left-3 top-2.5 w-4 h-4 text-gray-400 pointer-events-none"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                  <input
                    type="password"
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="Unesite lozinku"
                    required
                    className="w-full px-10 py-2 bg-[#fdf4ec] border border-orange-100 rounded-lg focus:ring-2 focus:ring-[#f4a261] focus:border-[#f4a261] text-[13px] text-gray-900"
                  />
                </div>
              </div>

              {/* Submit */}
              <button
                type="submit"
                disabled={loading}
                className="w-full bg-gradient-to-r from-[#e76f51] to-[#f4a261] hover:from-[#e85d44] hover:to-[#f1a550] text-white font-medium py-2.5 rounded-lg text-sm transition-all duration-200 disabled:opacity-50"
              >
                {loading ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-3.5 w-3.5 border-2 border-white border-t-transparent mr-2" />
                    Prijavljivanje...
                  </div>
                ) : (
                  'Prijavi se'
                )}
              </button>
            </form>

            {/* Links */}
            <div className="mt-2.5 text-center text-[11px]">
              <span className="text-gray-600">Nemate nalog? </span>
              <Link
                to="/register"
                className="text-[#e76f51] hover:text-[#d8583e] font-medium"
              >
                Registrujte se
              </Link>
            </div>

            {/* Test info */}
            <div className="mt-2.5 bg-[#fff5eb] rounded-lg p-2">
              <div className="text-center text-[10px] text-gray-700">
                <span className="font-medium mr-1">Test nalog:</span>
                <span className="bg-white rounded px-2 py-0.5 inline-block">
                  <strong>admin</strong> / <strong>password123</strong>
                </span>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
