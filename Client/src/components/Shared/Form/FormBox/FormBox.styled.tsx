import styled from "styled-components";

const FormBoxStyled = styled.div`
  background-color: ${(props) => props.theme.colors.background};
  margin: auto;
  padding: ${(props) => props.theme.paddings.large};
  width: 80%;
`;

export default FormBoxStyled;
