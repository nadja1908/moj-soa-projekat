import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Register = () => {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    role: 'tourist'
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Lozinke se ne poklapaju');
      return;
    }

    setLoading(true);
    const result = await register(formData);

    if (result.success) {
      navigate('/login');
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

      {/* Desni deo - kompaktniji card */}
      <div className="w-full md:w-1/2 flex items-center justify-center px-3">
        <div className="bg-white rounded-2xl shadow-md w-full max-w-sm border border-orange-100 p-4">
          <h2 className="text-center text-sm font-semibold text-[#222] mb-1">
            Registracija
          </h2>
          <p className="text-center text-[10px] text-[#777] mb-3">
            Započnite avanturu sa nama
          </p>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-2 py-1 rounded-md mb-2 text-[10px] text-center">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-2">
            {/* Uloge */}
            <div className="grid grid-cols-2 gap-2">
              {['tourist', 'guide'].map((r) => (
                <label
                  key={r}
                  className={`text-center py-1.5 text-[10px] rounded-md border cursor-pointer ${
                    formData.role === r
                      ? 'border-[#e76f51] bg-orange-50 text-[#e76f51]'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <input
                    type="radio"
                    name="role"
                    value={r}
                    checked={formData.role === r}
                    onChange={handleChange}
                    className="sr-only"
                  />
                  {r === 'tourist' ? 'Turista' : 'Vodič'}
                </label>
              ))}
            </div>

            {/* Ime i prezime */}
            <div className="grid grid-cols-2 gap-2">
              <input
                type="text"
                name="firstName"
                value={formData.firstName}
                onChange={handleChange}
                placeholder="Ime"
                required
                className="w-full px-2 py-1.5 text-[11px] bg-[#fdf4ec] border border-orange-100 rounded-md focus:ring-1 focus:ring-[#f4a261] focus:border-[#f4a261]"
              />
              <input
                type="text"
                name="lastName"
                value={formData.lastName}
                onChange={handleChange}
                placeholder="Prezime"
                required
                className="w-full px-2 py-1.5 text-[11px] bg-[#fdf4ec] border border-orange-100 rounded-md focus:ring-1 focus:ring-[#f4a261] focus:border-[#f4a261]"
              />
            </div>

            {/* Username */}
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              placeholder="Korisničko ime"
              required
              className="w-full px-2 py-1.5 text-[11px] bg-[#fdf4ec] border border-orange-100 rounded-md focus:ring-1 focus:ring-[#f4a261] focus:border-[#f4a261]"
            />

            {/* Email */}
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="Email"
              required
              className="w-full px-2 py-1.5 text-[11px] bg-[#fdf4ec] border border-orange-100 rounded-md focus:ring-1 focus:ring-[#f4a261] focus:border-[#f4a261]"
            />

            {/* Lozinke */}
            <div className="grid grid-cols-2 gap-2">
              <input
                type="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="Lozinka"
                required
                className="w-full px-2 py-1.5 text-[11px] bg-[#fdf4ec] border border-orange-100 rounded-md focus:ring-1 focus:ring-[#f4a261] focus:border-[#f4a261]"
              />
              <input
                type="password"
                name="confirmPassword"
                value={formData.confirmPassword}
                onChange={handleChange}
                placeholder="Ponovo"
                required
                className="w-full px-2 py-1.5 text-[11px] bg-[#fdf4ec] border border-orange-100 rounded-md focus:ring-1 focus:ring-[#f4a261] focus:border-[#f4a261]"
              />
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-gradient-to-r from-[#e76f51] to-[#f4a261] hover:from-[#e85d44] hover:to-[#f1a550] text-white font-medium py-1.5 rounded-md text-[12px] transition-all duration-200 disabled:opacity-50"
            >
              {loading ? 'Registracija...' : 'Registruj se'}
            </button>
          </form>

          <div className="mt-2 text-center text-[10px]">
            Već imate nalog?{' '}
            <Link
              to="/login"
              className="text-[#e76f51] hover:text-[#d8583e] font-medium"
            >
              Prijavite se
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Register;
