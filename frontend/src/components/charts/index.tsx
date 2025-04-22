import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  RadialLinearScale,
  Filler,
  ChartOptions,
} from 'chart.js';
import { Line, Radar } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  RadialLinearScale,
  Title,
  Tooltip,
  Legend,
  Filler
);

const months = ['Jan', 'Feb', 'Mar', 'Apr'];
const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

// Helper function to generate random data
const generateRandomData = (min: number, max: number, count: number) => {
  return Array.from({ length: count }, () => 
    Math.floor(Math.random() * (max - min + 1)) + min
  );
};

export const BookingTrendChart: React.FC = () => {
  const options: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: 'white',
        titleColor: 'black',
        bodyColor: 'black',
        borderColor: 'rgba(0,0,0,0.1)',
        borderWidth: 1,
        padding: 10,
        displayColors: false,
      },
    },
    scales: {
      x: {
        grid: {
          display: false,
        },
        ticks: {
          color: '#666',
        },
      },
      y: {
        grid: {
          color: 'rgba(0, 0, 0, 0.05)',
        },
        min: 0,
        max: 25,
        ticks: {
          stepSize: 5,
          color: '#666',
        },
      },
    },
    elements: {
      line: {
        tension: 0.4,
      },
      point: {
        radius: 4,
        hoverRadius: 6,
      },
    },
  };

  const data = {
    labels: months,
    datasets: [
      {
        data: generateRandomData(5, 20, 4),
        borderColor: '#6366F1',
        backgroundColor: '#6366F1',
        fill: false,
      },
    ],
  };

  return <Line options={options} data={data} />;
};

export const PlatformBookingChart: React.FC = () => {
  const options: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
        labels: {
          usePointStyle: true,
          pointStyle: 'circle',
          padding: 20,
          color: '#666',
        },
      },
      tooltip: {
        backgroundColor: 'white',
        titleColor: 'black',
        bodyColor: 'black',
        borderColor: 'rgba(0,0,0,0.1)',
        borderWidth: 1,
        padding: 10,
      },
    },
    scales: {
      x: {
        display: false,
      },
      y: {
        display: false,
      },
    },
    elements: {
      line: {
        tension: 0.4,
      },
      point: {
        radius: 0,
      },
    },
  };

  const data = {
    labels: Array.from({ length: 12 }, (_, i) => `Point ${i + 1}`),
    datasets: [
      {
        label: 'Booking.com',
        data: generateRandomData(60, 100, 12),
        borderColor: '#6366F1',
        backgroundColor: '#6366F1',
        tension: 0.4,
      },
      {
        label: 'Expedia',
        data: generateRandomData(40, 80, 12),
        borderColor: '#FFA500',
        backgroundColor: '#FFA500',
        tension: 0.4,
      },
      {
        label: 'Hotels.com',
        data: generateRandomData(20, 60, 12),
        borderColor: '#FF6B6B',
        backgroundColor: '#FF6B6B',
        tension: 0.4,
      },
    ],
  };

  return <Line options={options} data={data} />;
};

export const WeeklyVisitorsChart: React.FC = () => {
  const options: ChartOptions<'radar'> = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      r: {
        beginAtZero: true,
        grid: {
          color: 'rgba(0, 0, 0, 0.1)',
        },
        ticks: {
          display: false,
          maxTicksLimit: 5,
        },
        pointLabels: {
          color: '#666',
          font: {
            size: 10,
          },
        },
      },
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: 'white',
        titleColor: 'black',
        bodyColor: 'black',
        borderColor: 'rgba(0,0,0,0.1)',
        borderWidth: 1,
        padding: 10,
      },
    },
  };

  const data = {
    labels: days,
    datasets: [
      {
        data: generateRandomData(60, 90, 7),
        backgroundColor: 'rgba(99, 102, 241, 0.2)',
        borderColor: '#6366F1',
        borderWidth: 2,
        pointBackgroundColor: '#6366F1',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: '#6366F1',
      },
    ],
  };

  return <Radar options={options} data={data} />;
}; 