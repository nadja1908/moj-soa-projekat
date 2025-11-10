import React from 'react';

const TailwindTest = () => {
  return (
    <div className="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl">
      <div className="md:flex">
        <div className="md:shrink-0">
          <div className="h-48 w-full bg-gradient-to-r from-blue-500 to-purple-600 md:h-full md:w-48 flex items-center justify-center">
            <span className="text-white text-2xl font-bold">Tailwind CSS</span>
          </div>
        </div>
        <div className="p-8">
          <div className="uppercase tracking-wide text-sm text-indigo-500 font-semibold">
            Test Komponenta
          </div>
          <h2 className="block mt-1 text-lg leading-tight font-medium text-black">
            Tailwind CSS uspešno instaliran!
          </h2>
          <p className="mt-2 text-slate-500">
            Ova komponenta koristi čiste Tailwind CSS klase za stilizovanje.
          </p>
          <div className="mt-4">
            <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition-colors duration-200">
              Test Button
            </button>
            <button className="ml-2 bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition-colors duration-200">
              Drugi Button
            </button>
          </div>
          <div className="mt-4 flex space-x-2">
            <span className="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700">
              #tailwind
            </span>
            <span className="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700">
              #css
            </span>
            <span className="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700">
              #react
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TailwindTest;