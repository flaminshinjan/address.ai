import React from 'react';
import { Typography, Box, Breadcrumbs, Link as MuiLink } from '@mui/material';
import { Link } from 'react-router-dom';

interface PageHeaderProps {
  title: string;
  breadcrumbs?: Array<{
    text: string;
    href?: string;
  }>;
  action?: React.ReactNode;
}

export const PageHeader: React.FC<PageHeaderProps> = ({ title, breadcrumbs, action }) => {
  return (
    <Box mb={4}>
      {breadcrumbs && (
        <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
          <MuiLink component={Link} to="/" color="inherit">
            Home
          </MuiLink>
          {breadcrumbs.map((crumb, index) => (
            crumb.href ? (
              <MuiLink
                key={index}
                component={Link}
                to={crumb.href}
                color="inherit"
              >
                {crumb.text}
              </MuiLink>
            ) : (
              <Typography key={index} color="text.primary">
                {crumb.text}
              </Typography>
            )
          ))}
        </Breadcrumbs>
      )}
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Typography variant="h4" component="h1" gutterBottom={!action}>
          {title}
        </Typography>
        {action && <Box>{action}</Box>}
      </Box>
    </Box>
  );
}; 