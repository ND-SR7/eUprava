import styled from "styled-components";

const InputStyled = styled.input`
  border: 2px black solid;
  border-radius: ${(props) => props.theme.borderRadius.small};
  padding: 2px;
  margin: ${(props) => props.theme.margins.standard};
`;

export default InputStyled;
