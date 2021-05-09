import React from 'react';

type ResponseProps = {
  value: string;
};

const Response: React.FC<ResponseProps> = ({ value }) => {
  return (
    <div>
      <h5>{value}</h5>
    </div>
  );
};

export default Response;
