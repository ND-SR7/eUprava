import React, { useState } from 'react';
import { getMostPopularBrands } from '../../../services/StatisticsService';
import { Container, Table, Loader, ErrorMessage } from './MostPopularBrands.styled';
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';
import Button from '../../Shared/Button/Button';
import Input from '../../Shared/Input/Input';

interface BrandCount {
  brand: string;
  count: number;
}

const MostPopularBrands: React.FC = () => {
  const [year, setYear] = useState('');
  const [brands, setBrands] = useState<BrandCount[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [searchYear, setSearchYear] = useState<string | null>(null);

  const handleSearch = async () => {
    if (!/^\d{4}$/.test(year)) {
      setError('Please enter a valid year');
      setBrands([]);
      setSearchYear(null);
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      const data = await getMostPopularBrands(year);
      setBrands(data);
      setSearchYear(year);
      setError(null);
      if (data.length === 0) {
        setError(`No registered vehicles found for the year ${year}`);
      }
    } catch (err) {
      setError('Failed to fetch most popular brands');
      setBrands([]);
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

        pdf.save(`Most_Popular_Brands_Report_${searchYear}.pdf`);
      });
    }
  };

  const isDownloadDisabled = brands.length === 0;

  return (
    <Container>
      <h1>Search Most Popular Vehicle Brands by Year</h1>
      <div>
        <Input
          type="text"
          label="Year"
          id="year"
          attrName="year"
          handleChange={(e) => setYear(e.target.value)}
          data={year}
        />
        <Button 
          label="Search" 
          buttonType="button" 
          onClick={handleSearch} 
        />
      </div>
      {loading && <Loader>Loading...</Loader>}
      {error && searchYear && <ErrorMessage>{error}</ErrorMessage>}
      {brands.length > 0 && searchYear !== null && (
        <>
          <div id="report-content">
            <Table>
              <thead>
                <tr>
                  <th>Brand</th>
                  <th>Number of Registered Vehicles</th>
                </tr>
              </thead>
              <tbody>
                {brands.map((brand) => (
                  <tr key={brand.brand}>
                    <td>{brand.brand}</td>
                    <td>{brand.count}</td>
                  </tr>
                ))}
              </tbody>
            </Table>
          </div>
          <Button 
            label="Download PDF" 
            buttonType="button" 
            onClick={generatePDF} 
            disabled={isDownloadDisabled} 
          />
        </>
      )}
    </Container>
  );
};

export default MostPopularBrands;
