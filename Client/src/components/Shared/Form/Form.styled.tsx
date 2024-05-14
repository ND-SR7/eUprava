import styled from 'styled-components';

const FormStyled = styled.article`
  border: ${(props) => props.theme.borders.standardOrange};
  padding: ${(props) => props.theme.paddings.standard};
`;

export default FormStyled;
