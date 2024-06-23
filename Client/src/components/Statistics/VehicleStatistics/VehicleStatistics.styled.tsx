import styled from 'styled-components';

export const StyledButton = styled.button`
  background-color: #4CAF50;
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  margin: 4px 2px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.3s;

  &:hover {
    background-color: #45a049;
  }
`;

export const VehicleStatisticsTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
  font-size: 1em;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  min-width: 400px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);

  th, td {
    padding: 12px 15px;
    border: 1px solid #ddd;
  }

  th {
    background-color: #009879;
    color: #ffffff;
    text-align: left;
    font-weight: bold;
  }

  tr {
    border-bottom: 1px solid #dddddd;
  }

  tr:nth-of-type(even) {
    background-color: #f3f3f3;
  }

  tr:last-of-type {
    border-bottom: 2px solid #009879;
  }

  tr:hover {
    background-color: #f1f1f1;
  }
`;

export const NoVehicleStatisticsMessage = styled.div`
  text-align: center;
  margin: 20px 0;
  color: #777;
`;

export const ModalContent = styled.div`
  max-height: 80vh;
  overflow-y: auto;
  padding: 20px;
`;

export const ModalHeader = styled.div`
  margin-bottom: 20px;
  text-align: center;
`;
