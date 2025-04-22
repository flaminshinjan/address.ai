import { Theme, ThemeOptions } from '@mui/material/styles';
import { DataGridProps } from '@mui/x-data-grid';

declare module '@mui/material/styles' {
  interface ComponentNameToClassKey {
    MuiDataGrid: DataGridProps;
  }

  interface Components<Theme = unknown> {
    MuiDataGrid?: {
      defaultProps?: Partial<DataGridProps>;
      styleOverrides?: {
        root?: React.CSSProperties;
        cell?: React.CSSProperties;
        columnHeader?: React.CSSProperties;
        row?: {
          '&:nth-of-type(odd)'?: React.CSSProperties;
          '&:hover'?: React.CSSProperties;
        };
      };
    };
  }
} 