import styled from "styled-components";

const SelectStyled = styled.select`
  margin: 10px 0;
  padding: 10px;
  width: 100%;
  max-width: 300px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  background-color: ${(props) => props.theme.colors.inputBackground};
  color: ${(props) => props.theme.colors.text};
`;

export default SelectStyled;
