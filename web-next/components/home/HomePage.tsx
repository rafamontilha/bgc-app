'use client';

/**
 * Home Page Component
 * Main landing page orchestrating all sections
 * Aligned with original globalconnect-bridge structure
 */

import React, { useCallback } from 'react';
import { Box } from '@mui/material';
import { Header } from './Header';
import { HeroSection } from './HeroSection';
import { SimulatorSection } from './SimulatorSection';
import { FeaturesSection } from './FeaturesSection';
import { CTASection } from './CTASection';
import { Footer } from './Footer';

export function HomePage(): React.ReactElement {
  const handleLoginClick = useCallback(() => {
    // Navigate to login page
    window.location.href = '/login';
  }, []);

  const handleGetStarted = useCallback(() => {
    const simulatorSection = document.getElementById('simulator');
    if (simulatorSection) {
      simulatorSection.scrollIntoView({ behavior: 'smooth' });
    }
  }, []);

  const handleLearnMore = useCallback(() => {
    const featuresSection = document.getElementById('features');
    if (featuresSection) {
      featuresSection.scrollIntoView({ behavior: 'smooth' });
    }
  }, []);

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      {/* Header */}
      <Header onLoginClick={handleLoginClick} />

      {/* Main content */}
      <Box component="main" sx={{ flex: 1 }}>
        {/* Hero Section - Value Proposition */}
        <HeroSection
          onGetStartedClick={handleGetStarted}
          onLearnMoreClick={handleLearnMore}
        />

        {/* Simulator Section - After scroll */}
        <SimulatorSection />

        {/* Features Section - 3 main benefits */}
        <FeaturesSection />

        {/* CTA Section - Final call-to-action */}
        <CTASection />
      </Box>

      {/* Footer */}
      <Footer />
    </Box>
  );
}
