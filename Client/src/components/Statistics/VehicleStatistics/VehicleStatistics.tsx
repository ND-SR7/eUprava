import React, { useState, useEffect } from 'react';
import { getVehicleStatisticsByYear } from '../../../services/StatisticsService';
import { VehicleStatisticsTable, NoVehicleStatisticsMessage, ModalContent, ModalHeader } from './VehicleStatistics.styled';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const VehicleStatistics: React.FC = () => {
  const [vehicleStatistics, setVehicleStatistics] = useState<{ [year: number]: number }>({});

  useEffect(() => {
    getVehicleStatisticsByYear()
      .then((data) => {
        setVehicleStatistics(data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  return (
    <ModalContent>
      <ModalHeader>
        <h2>Vehicle Statistics by Year</h2>
      </ModalHeader>
      {Object.keys(vehicleStatistics).length > 0 ? (
        <>
          <VehicleStatisticsTable>
            <thead>
              <tr>
                <th>Year</th>
                <th>Count</th>
              </tr>
            </thead>
            <tbody>
              {Object.entries(vehicleStatistics).map(([year, count]) => (
                <tr key={year}>
                  <td>{year}</td>
                  <td>{count}</td>
                </tr>
              ))}
            </tbody>
          </VehicleStatisticsTable>
          <div style={{ width: '100%', height: 300 }}>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart
                data={Object.entries(vehicleStatistics).map(([year, count]) => ({ year, count }))}
                margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="year" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="count" stroke="#8884d8" activeDot={{ r: 8 }} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </>
      ) : (
        <NoVehicleStatisticsMessage>No statistics available</NoVehicleStatisticsMessage>
      )}
    </ModalContent>
  );
};

export default VehicleStatistics;
