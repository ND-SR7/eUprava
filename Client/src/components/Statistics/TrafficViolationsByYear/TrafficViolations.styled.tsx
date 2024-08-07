import styled from 'styled-components';

export const Container = styled.div`
    padding: 20px;
    max-width: 800px;
    margin: 0 auto;
`;

export const Title = styled.h1`
    text-align: center;
    margin-bottom: 20px;
`;

export const Form = styled.form`
    display: flex;
    justify-content: center;
    margin-bottom: 20px;
`;

export const Table = styled.table`
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
`;

export const TableHeader = styled.th`
    background-color: #009879;
    color: white;
    padding: 10px;
    border: 1px solid #dee2e6;
    text-align: left;
`;

export const TableRow = styled.tr`
    &:nth-child(even) {
        background-color: #f2f2f2;
    }
    &:last-child {
        font-weight: bold;
    }
`;

export const TableData = styled.td`
    padding: 10px;
    border: 1px solid #dee2e6;
`;

export const Loader = styled.div`
    text-align: center;
    margin-top: 20px;
`;

export const ErrorMessage = styled.div`
    color: red;
    text-align: center;
    margin-top: 20px;
`;
