import React, { useState } from 'react';
import { getRegisteredVehiclesByYear } from '../../../services/StatisticsService';
import { Container, Input, Button, Table, Message, Loader, ErrorMessage } from './RegisteredVehicles.styled';

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
      )}
    </Container>
  );
};

export default RegisteredVehicles;
