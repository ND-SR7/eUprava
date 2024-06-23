import styled from 'styled-components';

export const Container = styled.div`
  text-align: center;
  margin: 20px;
`;

export const Input = styled.input`
  padding: 10px;
  margin-right: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
`;

export const Button = styled.button`
  padding: 10px;
  font-size: 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;

  &:hover {
    background-color: #45a049;
  }
`;

export const Table = styled.table`
  margin: 20px auto;
  border-collapse: collapse;
  width: 50%;

  th, td {
    padding: 10px;
    border: 1px solid #ccc;
    text-align: center;
  }

  th {
    background-color: #f4f4f4;
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
