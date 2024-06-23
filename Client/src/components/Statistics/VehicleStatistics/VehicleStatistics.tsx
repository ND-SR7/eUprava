import React, { useState, useEffect } from 'react';
import { getVehicleStatisticsByYear } from '../../../services/StatisticsService';
import { VehicleStatisticsTable, NoVehicleStatisticsMessage, ModalContent, ModalHeader, StyledButton, DownloadButton } from './VehicleStatistics.styled';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';

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

  const generatePDF = () => {
    const input = document.getElementById('vehicle-statistics-content');
    if (input) {
      html2canvas(input).then((canvas) => {
        const imgData = canvas.toDataURL('image/png');
        const pdf = new jsPDF();
        const imgWidth = 190;
        const pageHeight = 295;
        const imgHeight = (canvas.height * imgWidth) / canvas.width;
        let heightLeft = imgHeight;

        let position = 10;

        pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
        heightLeft -= pageHeight;

        while (heightLeft >= 0) {
          position = heightLeft - imgHeight;
          pdf.addPage();
          pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
          heightLeft -= pageHeight;
        }

        pdf.save('Vehicle_Statistics_Report.pdf');
      });
    }
  };

  return (
    <ModalContent>
      <div id="vehicle-statistics-content">
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
      </div>
      <DownloadButton onClick={generatePDF}>Download PDF</DownloadButton>
    </ModalContent>
  );
};

export default VehicleStatistics;
