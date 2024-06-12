import React from 'react'
import GetAllTrafficViolations from '../../components/Police/TrafficViolation/GetAllTrafficViolation'
import Button from '../../components/Shared/Button/Button'

const GetAllPages = () => {
  return (
    <>
        <GetAllTrafficViolations></GetAllTrafficViolations>
        <br />
        <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
        <br />
    </>
  
  );
};

export default GetAllPages;
