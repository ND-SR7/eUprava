import styled from 'styled-components';

export const Container = styled.div`
  text-align: center;
  margin: 20px;
`;

export const Table = styled.table`
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

export const Message = styled.p`
  font-size: 18px;
  margin-top: 20px;
`;

export const ErrorMessage = styled.p`
  font-size: 18px;
  margin-top: 20px;
  color: red;
`;

export const Loader = styled.div`
  font-size: 18px;
  margin-top: 20px;
`;
