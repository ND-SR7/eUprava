import React, { useState } from 'react';
import { getRegisteredVehiclesByYear } from '../../../services/StatisticsService';
import { Container, Input, Button, Table, Loader, ErrorMessage, DownloadButton } from './RegisteredVehicles.styled';
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';

const RegisteredVehicles: React.FC = () => {
  const [year, setYear] = useState('');
  const [vehicleCount, setVehicleCount] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [searchYear, setSearchYear] = useState<string | null>(null);

  const handleSearch = async () => {
    if (!/^\d{4}$/.test(year)) {
      setError('Please enter a valid year');
      setVehicleCount(null);
      return;
    }
    setLoading(true);
    try {
      const data = await getRegisteredVehiclesByYear(year);
      setVehicleCount(data.count);
      setSearchYear(year);
      setError(null);
    } catch (err) {
      setError('Failed to fetch registered vehicles');
      setVehicleCount(null);
    }
    setLoading(false);
  };

  const generatePDF = () => {
    const input = document.getElementById('report-content');
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

        pdf.save(`Registered_Vehicles_Report_${searchYear}.pdf`);
      });
    }
  };

  return (
    <Container>
      <h1>Search Registered Vehicles by Year</h1>
      <div>
        <Input
          type="text"
          value={year}
          onChange={(e) => setYear(e.target.value)}
          placeholder="Enter year"
        />
        <Button onClick={handleSearch}>Search</Button>
      </div>
      {loading && <Loader>Loading...</Loader>}
      {error && <ErrorMessage>{error}</ErrorMessage>}
      {vehicleCount !== null && searchYear !== null && (
        <>
          <div id="report-content">
            <Table>
              <thead>
                <tr>
                  <th>Year</th>
                  <th>Number of Registered Vehicles</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>{searchYear}</td>
                  <td>{vehicleCount}</td>
                </tr>
              </tbody>
            </Table>
          </div>
          <DownloadButton onClick={generatePDF}>Download PDF</DownloadButton>
        </>
      )}
    </Container>
  );
};

export default RegisteredVehicles;
