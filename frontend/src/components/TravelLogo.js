// src/components/TravelLogo.jsx (ili gde već držiš)

import React from 'react';

const TravelLogo = ({ className = "w-8 h-8" }) => {
  return (
    <div className={`${className} relative flex items-center justify-center`}>
      <img 
        src="/logo.png" 
        alt="Travel Tourist Logo" 
        className="w-full h-full object-contain"
      />
    </div>
  );
};

export default TravelLogo;
